package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

// App struct
type App struct {
	ctx context.Context
}

// BrewPackage 定义一个结构体，方便前端使用
type BrewPackage struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Status  string `json:"status"` // "started", "stopped", 或 "none" (不是服务)
}

type BrewData struct {
	Formulae []BrewPackage `json:"formulae"`
	Casks    []BrewPackage `json:"casks"`
	Taps     []BrewTap     `json:"taps"` // 新增：tap 列表
}
type ServiceInfo struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

// BrewTap 定义 Homebrew tap 的数据结构
type BrewTap struct {
	Name        string `json:"name"`        // tap 名称，如 "homebrew/core"
	Official    bool   `json:"official"`    // 是否为官方 tap
	URL         string `json:"url"`         // tap 的 Git URL
	Description string `json:"description"` // tap 描述（可选）
}

// ActionResponse 操作结果返回
type ActionResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 获取 brew 的路径
func getBrewPath() string {
	// 检查 M1/M2 路径
	if _, err := os.Stat("/opt/homebrew/bin/brew"); err == nil {
		return "/opt/homebrew/bin/brew"
	}
	// 检查 Intel 路径
	if _, err := os.Stat("/usr/local/bin/brew"); err == nil {
		return "/usr/local/bin/brew"
	}
	return "brew" // 保底方案
}

// GetBrewData 获取 brew 数据
func (a *App) GetBrewData() BrewData {
	// 1. 获取所有服务状态
	services := make(map[string]string)
	serviceRaw, _ := exec.Command(getBrewPath(), "services", "info", "--all", "--json").Output()
	var serviceList []ServiceInfo
	json.Unmarshal(serviceRaw, &serviceList)
	for _, s := range serviceList {
		services[s.Name] = s.Status
	}

	// 2. 获取 taps 数据
	taps := a.GetBrewTaps()

	// 3. 获取并包装数据
	return BrewData{
		Formulae: fetchWithStatus("--formula", services),
		Casks:    fetchWithStatus("--cask", services),
		Taps:     taps,
	}
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// 重点：为 GUI 进程手动注入 PATH，否则执行 brew 命令会卡死或直接报错找不到命令
	path := os.Getenv("PATH")
	// 兼容 M1/M2/M3 和 Intel Mac，同时保留原有路径
	newPath := "/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin:" + path
	os.Setenv("PATH", newPath)
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// brew list 命令并返回字符串数组
func fetchWithStatus(flag string, serviceMap map[string]string) []BrewPackage {
	out, _ := exec.Command(getBrewPath(), "list", "--versions", flag).Output()
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var packages []BrewPackage

	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			name := parts[0]
			version := strings.Join(parts[1:], " ")
			status, isService := serviceMap[name]
			if !isService {
				status = "none_tool" // 普通工具，没有服务状态，服务停止会显示 none
			}

			packages = append(packages, BrewPackage{
				Name:    name,
				Version: version,
				Status:  status,
			})
		}
	}
	// --- 排序逻辑开始 ---
	sort.Slice(packages, func(i, j int) bool {
		// 定义状态权重
		priority := map[string]int{
			"started":   1,
			"stopped":   2,
			"none":      3,
			"none_tool": 4,
		}

		// 如果状态不同，按优先级排（1最小，排最前）
		if packages[i].Status != packages[j].Status {
			return priority[packages[i].Status] < priority[packages[j].Status]
		}
		// 如果状态相同，按名称字母顺序排
		return packages[i].Name < packages[j].Name
	})
	// --- 排序逻辑结束 ---
	return packages
}

// StartService 启动指定的 Brew 服务
func (a *App) StartService(name string) ActionResponse {
	// 执行命令: brew services start <name>
	out, err := exec.Command(getBrewPath(), "services", "start", name).CombinedOutput()
	if err != nil {
		// 这样即使失败，也会把错误信息返回给前端的 alert
		return ActionResponse{
			Success: false,
			Message: fmt.Sprintf("后端执行失败: %v, 输出: %s", err, string(out)),
		}
	}
	return ActionResponse{Success: true, Message: "服务 " + name + "已启动"}
}

// StopService 停止指定的 Brew 服务
func (a *App) StopService(name string) ActionResponse {
	// 执行命令: brew services stop <name>
	out, err := exec.Command(getBrewPath(), "services", "stop", name).CombinedOutput()
	if err != nil {
		return ActionResponse{
			Success: false,
			Message: fmt.Sprintf("停止失败(%v): %s", err, string(out)),
		}
		// return ActionResponse{Success: false, Message: "停止失败：" + err.Error()}
	}
	return ActionResponse{Success: true, Message: "服务 " + name + "已停止"}
}

// GetAppIcon 获取应用的图标并转为 Base64
func (a *App) GetAppIcon(appName string) string {
	// 1. 建立手动映射表 (解决包名和 App 名完全对不上的情况)
	mismatchedNames := map[string]string{
		"iterm2":             "iTerm",
		"docker-desktop":     "Docker",
		"google-chrome":      "Google Chrome",
		"visual-studio-code": "Visual Studio Code",
		"dbeaver-community":  "DBeaver", // 针对 DBeaver 的映射
	}

	// 2. 生成可能的 App 搜索名称
	searchNames := []string{appName}
	if realName, ok := mismatchedNames[appName]; ok {
		searchNames = append(searchNames, realName)
	}

	// 补充常规变形：首字母大写、全大写
	searchNames = append(searchNames, capitalize(appName))
	searchNames = append(searchNames, strings.ToUpper(appName))

	var appPath string
	// --- 第一阶段：尝试所有已知可能的精确路径 ---
	for _, name := range searchNames {
		path := filepath.Join("/Applications", name+".app")
		if _, err := os.Stat(path); err == nil {
			appPath = path
			break
		}
	}

	// --- 第二阶段：如果还没找到，才进行谨慎的模糊匹配 ---
	if appPath == "" {
		// 如果还是找不到，尝试在 /Applications 目录下遍历一遍，看看有没有包含关键字的
		files, _ := os.ReadDir("/Applications")
		for _, f := range files {
			// 确保找到的是 .app 结尾的文件夹，且排除干扰项
			fileName := f.Name()
			if strings.HasSuffix(fileName, ".app") &&
				strings.Contains(strings.ToLower(fileName), strings.ToLower(appName)) {
				appPath = filepath.Join("/Applications", fileName)
				break
			}
		}
	}

	if appPath == "" {
		return ""
	}

	// 3. 读取 Info.plist 定位图标 (iTerm2 这里的图标名通常是 iTerm2.icns)
	iconFileName := "AppIcon.icns" // 默认
	plistPath := filepath.Join(appPath, "Contents", "Info.plist")
	if data, err := os.ReadFile(plistPath); err == nil {
		content := string(data)
		if idx := strings.Index(content, "CFBundleIconFile"); idx != -1 {
			sub := content[idx:]
			start := strings.Index(sub, "<string>") + 8
			end := strings.Index(sub, "</string>")
			if start > 7 && end > start {
				iconFileName = sub[start:end]
				if !strings.HasSuffix(iconFileName, ".icns") {
					iconFileName += ".icns"
				}
			}
		}
	}

	iconPath := filepath.Join(appPath, "Contents", "Resources", iconFileName)
	// 兜底方案：如果指定的图标不存在，扫描 Resources 下任何一个 .icns
	if _, err := os.Stat(iconPath); err != nil {
		resFiles, _ := os.ReadDir(filepath.Join(appPath, "Contents", "Resources"))
		for _, rf := range resFiles {
			if strings.HasSuffix(strings.ToLower(rf.Name()), ".icns") {
				iconPath = filepath.Join(appPath, "Contents", "Resources", rf.Name())
				break
			}
		}
	}

	// 4. 转换并返回 (带唯一 AppName 的临时文件防止图标错乱)
	tmpPng := filepath.Join(os.TempDir(), fmt.Sprintf("icon_%s.png", appName))
	// tmpPng := filepath.Join(os.TempDir(), fmt.Sprintf("brew_icon_%s.png", appName))
	defer os.Remove(tmpPng)

	cmd := exec.Command("sips", "-s", "format", "png", iconPath, "--out", tmpPng)
	if err := cmd.Run(); err != nil {
		return ""
	}

	imgData, _ := os.ReadFile(tmpPng)
	return base64.StdEncoding.EncodeToString(imgData)
}

func capitalize(s string) string {
	if s == "" {
		return ""
	}
	// 将字符串转为 rune 切片以正确处理 Unicode 字符
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// RestartService 重启服务
func (a *App) RestartService(name string) ActionResponse {
	cmd := exec.Command(getBrewPath(), "services", "restart", name)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return ActionResponse{
			Success: false,
			Message: fmt.Sprintf("重启失败: %s", string(out)),
		}
	}
	return ActionResponse{
		Success: true,
		Message: fmt.Sprintf("服务 %s 已重启", name),
	}
}

// ------------------------ Tap 管理功能 ------------------------

// GetBrewTaps 获取所有已安装的 taps
func (a *App) GetBrewTaps() []BrewTap {
	brewPath := getBrewPath()
	// 执行 brew tap 命令获取已安装的 tap 列表
	cmd := exec.Command(brewPath, "tap")
	output, err := cmd.Output()
	if err != nil {
		return []BrewTap{}
	}

	var taps []BrewTap
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		tap := parseTapInfo(line, brewPath)
		taps = append(taps, tap)
	}

	return taps
}

// parseTapInfo 解析单个 tap 的信息
func parseTapInfo(tapName, brewPath string) BrewTap {
	tap := BrewTap{
		Name: tapName,
	}

	// 检查是否为官方 tap
	officialTaps := []string{
		"homebrew/core",
		"homebrew/cask",
		"homebrew/cask-fonts",
		"homebrew/cask-drivers",
		"homebrew/cask-versions",
		"homebrew/bundle",
		"homebrew/services",
	}

	for _, official := range officialTaps {
		if tapName == official {
			tap.Official = true
			break
		}
	}

	// 获取 tap 的 Git URL
	cmd := exec.Command(brewPath, "tap-info", tapName, "--json")
	output, err := cmd.Output()
	if err == nil {
		// 解析 JSON 获取 tap 信息
		var tapInfos []struct {
			Name        string `json:"name"`
			Remote      string `json:"remote"`
			Description string `json:"desc"`
		}
		if json.Unmarshal(output, &tapInfos) == nil && len(tapInfos) > 0 {
			tap.URL = tapInfos[0].Remote
			tap.Description = tapInfos[0].Description
		}
	}

	return tap
}

// AddTap 添加新的 tap
func (a *App) AddTap(tapName string) ActionResponse {
	brewPath := getBrewPath()

	// 检查 tap 是否已经存在
	taps := a.GetBrewTaps()
	for _, tap := range taps {
		if tap.Name == tapName {
			return ActionResponse{
				Success: false,
				Message: fmt.Sprintf("Tap '%s' 已经存在", tapName),
			}
		}
	}

	// 执行 brew tap 命令添加 tap
	cmd := exec.Command(brewPath, "tap", tapName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ActionResponse{
			Success: false,
			Message: fmt.Sprintf("添加 tap 失败: %s", string(output)),
		}
	}

	return ActionResponse{
		Success: true,
		Message: fmt.Sprintf("Tap '%s' 添加成功", tapName),
	}
}

// RemoveTap 移除 tap
func (a *App) RemoveTap(tapName string, force bool) ActionResponse {
	brewPath := getBrewPath()

	// 检查是否为官方 tap，警告用户
	tap := parseTapInfo(tapName, brewPath)
	if tap.Official && !force {
		return ActionResponse{
			Success: false,
			Message: fmt.Sprintf("警告: '%s' 是官方 tap，不建议移除。如确认移除，请使用强制删除。", tapName),
		}
	}

	// 执行 brew untap 命令移除 tap
	cmd := exec.Command(brewPath, "untap", tapName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ActionResponse{
			Success: false,
			Message: fmt.Sprintf("移除 tap 失败: %s", string(output)),
		}
	}

	return ActionResponse{
		Success: true,
		Message: fmt.Sprintf("Tap '%s' 移除成功", tapName),
	}
}

// UpdateTap 更新指定 tap
func (a *App) UpdateTap(tapName string) ActionResponse {
	brewPath := getBrewPath()

	// 检查 tap 是否存在
	taps := a.GetBrewTaps()
	found := false
	for _, tap := range taps {
		if tap.Name == tapName {
			found = true
			break
		}
	}

	if !found {
		return ActionResponse{
			Success: false,
			Message: fmt.Sprintf("Tap '%s' 不存在", tapName),
		}
	}

	// 执行 brew tap 命令更新 tap
	cmd := exec.Command(brewPath, "tap", "--repair", tapName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ActionResponse{
			Success: false,
			Message: fmt.Sprintf("更新 tap 失败: %s", string(output)),
		}
	}

	return ActionResponse{
		Success: true,
		Message: fmt.Sprintf("Tap '%s' 更新成功", tapName),
	}
}

// UpdateAllTaps 更新所有 taps
func (a *App) UpdateAllTaps() ActionResponse {
	brewPath := getBrewPath()

	// 执行 brew update 命令更新所有 taps
	cmd := exec.Command(brewPath, "update")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ActionResponse{
			Success: false,
			Message: fmt.Sprintf("更新所有 taps 失败: %s", string(output)),
		}
	}

	return ActionResponse{
		Success: true,
		Message: "所有 taps 更新成功",
	}
}

// GetTapPackageCount 获取 tap 中的 package 数量（formulae + casks）
func (a *App) GetTapPackageCount(tapName string) map[string]int {
	result := map[string]int{
		"formulae": 0,
		"casks":    0,
	}

	brewPath := getBrewPath()

	// 获取 tap 中的 formulae
	cmd := exec.Command(brewPath, "info", "--json=v2", tapName)
	output, err := cmd.Output()
	if err == nil {
		var info struct {
			Casks    []interface{} `json:"casks"`
			Formulae []interface{} `json:"formulae"`
		}
		if json.Unmarshal(output, &info) == nil {
			result["formulae"] = len(info.Formulae)
			result["casks"] = len(info.Casks)
		}
	}

	return result
}

// SearchPackages 搜索可安装的软件包和 tap
func (a *App) SearchPackages(query string) map[string]interface{} {
	brewPath := getBrewPath()

	// 搜索 formulae
	formulaCmd := exec.Command(brewPath, "search", "--formula", query)
	formulaOutput, _ := formulaCmd.Output()

	// 搜索 casks
	caskCmd := exec.Command(brewPath, "search", "--cask", query)
	caskOutput, _ := caskCmd.Output()

	// 解析结果
	formulae := parseSearchResults(string(formulaOutput))
	casks := parseSearchResults(string(caskOutput))

	// 搜索 taps（通过 GitHub API 搜索 homebrew- 开头的仓库）
	taps := a.searchTaps(query)

	return map[string]interface{}{
		"formulae": formulae,
		"casks":    casks,
		"taps":     taps,
		"total":    len(formulae) + len(casks) + len(taps),
	}
}

// searchTaps 搜索可用的 taps
func (a *App) searchTaps(query string) []map[string]string {
	var taps []map[string]string

	// 使用 curl 搜索 GitHub 上 homebrew- 开头的仓库
	searchURL := fmt.Sprintf("https://api.github.com/search/repositories?q=%s+homebrew+in:name&sort=stars&order=desc&per_page=10", query)

	cmd := exec.Command("curl", "-s", "-H", "Accept: application/vnd.github.v3+json", searchURL)
	output, err := cmd.Output()
	if err != nil {
		return taps
	}

	var result struct {
		Items []struct {
			FullName    string `json:"full_name"`
			Description string `json:"description"`
			HTMLURL     string `json:"html_url"`
			Stargazers  int    `json:"stargazers_count"`
		} `json:"items"`
	}

	if json.Unmarshal(output, &result) == nil {
		for _, item := range result.Items {
			// 只返回包含 homebrew- 的仓库
			if strings.Contains(item.FullName, "homebrew-") {
				tapName := strings.Replace(item.FullName, "homebrew-", "", 1)
				tapName = strings.Replace(tapName, "/homebrew-", "/", 1)
				taps = append(taps, map[string]string{
					"name":        tapName,
					"fullName":    item.FullName,
					"description": item.Description,
					"url":         item.HTMLURL,
					"stars":       fmt.Sprintf("%d", item.Stargazers),
				})
			}
		}
	}

	return taps
}

// parseSearchResults 解析搜索结果
func parseSearchResults(output string) []string {
	var results []string
	lines := strings.Split(strings.TrimSpace(output), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		// 跳过空行和提示信息
		if line == "" || strings.Contains(line, "==>") || strings.Contains(line, " formulae") || strings.Contains(line, " casks") {
			continue
		}
		results = append(results, line)
	}

	return results
}

// InstallPackage 安装软件包
func (a *App) InstallPackage(name string, isCask bool) ActionResponse {
	brewPath := getBrewPath()

	// 先搜索确认包是否存在
	searchResult := a.SearchPackages(name)
	formulae := searchResult["formulae"].([]string)
	casks := searchResult["casks"].([]string)

	// 检查是否找到完全匹配或包含该名称的包
	found := false
	if !isCask {
		for _, f := range formulae {
			if f == name {
				found = true
				break
			}
		}
	} else {
		for _, c := range casks {
			if c == name {
				found = true
				break
			}
		}
	}

	if !found {
		return ActionResponse{
			Success: false,
			Message: fmt.Sprintf("未找到软件包 '%s'，请先搜索确认", name),
		}
	}

	// 执行安装命令
	var cmd *exec.Cmd
	if isCask {
		cmd = exec.Command(brewPath, "install", "--cask", name)
	} else {
		cmd = exec.Command(brewPath, "install", name)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return ActionResponse{
			Success: false,
			Message: fmt.Sprintf("安装失败: %s", string(output)),
		}
	}

	return ActionResponse{
		Success: true,
		Message: fmt.Sprintf("'%s' 安装成功", name),
	}
}

// ------------------------ Docker 容器管理功能 ------------------------

// DockerContainer Docker 容器信息
type DockerContainer struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Image  string `json:"image"`
	Status string `json:"status"` // "running", "exited", "paused", etc.
	State  string `json:"state"`  // 容器状态
	Ports  string `json:"ports"`  // 端口映射
}

// GetDockerContainers 获取所有 Docker 容器
func (a *App) GetDockerContainers() ActionResponse {
	// 检查 docker 命令是否存在
	if _, err := exec.LookPath("docker"); err != nil {
		return ActionResponse{Success: false, Message: "Docker 未安装"}
	}

	cmd := exec.Command("docker", "ps", "-a", "--format", "json")
	output, err := cmd.CombinedOutput()
	if err != nil {
		errMsg := string(output)
		if strings.Contains(errMsg, "permission denied") || strings.Contains(errMsg, "Is the docker daemon running") {
			return ActionResponse{Success: false, Message: "Docker 未启动或权限不足"}
		}
		return ActionResponse{Success: false, Message: fmt.Sprintf("Docker 错误: %s", errMsg)}
	}

	var containers []DockerContainer
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		var c struct {
			ID     string `json:"ID"`
			Names  string `json:"Names"`
			Image  string `json:"Image"`
			Status string `json:"Status"`
			State  string `json:"State"`
			Ports  string `json:"Ports"`
		}
		if err := json.Unmarshal([]byte(line), &c); err != nil {
			continue
		}
		containers = append(containers, DockerContainer{
			ID: c.ID, Name: c.Names, Image: c.Image,
			Status: c.Status, State: c.State, Ports: c.Ports,
		})
	}

	return ActionResponse{Success: true, Data: containers}
}

// shortID 安全截取容器 ID（防止越界 panic）
func shortID(id string) string {
	if len(id) > 12 {
		return id[:12]
	}
	return id
}

// StartDockerContainer 启动 Docker 容器
func (a *App) StartDockerContainer(containerID string) ActionResponse {
	cmd := exec.Command("docker", "start", containerID)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ActionResponse{
			Success: false,
			Message: fmt.Sprintf("启动失败: %s", string(output)),
		}
	}
	return ActionResponse{
		Success: true,
		Message: fmt.Sprintf("容器 %s 已启动", shortID(containerID)),
	}
}

// StopDockerContainer 停止 Docker 容器
func (a *App) StopDockerContainer(containerID string) ActionResponse {
	cmd := exec.Command("docker", "stop", containerID)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ActionResponse{
			Success: false,
			Message: fmt.Sprintf("停止失败: %s", string(output)),
		}
	}
	return ActionResponse{
		Success: true,
		Message: fmt.Sprintf("容器 %s 已停止", shortID(containerID)),
	}
}
