package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

// App struct
type App struct {
	ctx context.Context
}

// 定义一个结构体，方便前端使用
type BrewPackage struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Status  string `json:"status"` // "started", "stopped", 或 "none" (不是服务)
}

type BrewData struct {
	Formulae []BrewPackage `json:"formulae"`
	Casks    []BrewPackage `json:"casks"`
}
type ServiceInfo struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

// 操作结果返回
type ActionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

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

func (a *App) GetBrewData() BrewData {
	// 1. 获取所有服务状态
	services := make(map[string]string)
	serviceRaw, _ := exec.Command(getBrewPath(), "services", "info", "--all", "--json").Output()
	var serviceList []ServiceInfo
	json.Unmarshal(serviceRaw, &serviceList)
	for _, s := range serviceList {
		services[s.Name] = s.Status
	}

	// 2. 获取并包装数据
	return BrewData{
		Formulae: fetchWithStatus("--formula", services),
		Casks:    fetchWithStatus("--cask", services),
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
