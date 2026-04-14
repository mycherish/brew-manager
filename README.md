# 🍺 Brew Manager

> **基于 Wails + Vue 3 构建的 macOS 原生质感 Homebrew 管理工具**

[![Wails](https://img.shields.io/badge/Built%20with-Wails-red.svg)](https://wails.io/)
[![Vue](https://img.shields.io/badge/Frontend-Vue%203-brightgreen.svg)](https://vuejs.org/)
[![Go](https://img.shields.io/badge/Backend-Go-blue.svg)](https://go.dev/)
[![Platform](https://img.shields.io/badge/Platform-macOS%20(Apple%20Silicon)-lightgrey.svg)]()
[![Release](https://img.shields.io/github/v/release/mycherish/brew-manager)](https://github.com/mycherish/brew-manager/releases)
[![Stars](https://img.shields.io/github/stars/mycherish/brew-manager)](https://github.com/mycherish/brew-manager)

`Brew Manager` 是一个为 macOS 用户打造的轻量级 Homebrew 图形界面工具。它不仅能让你一眼看清系统安装的所有 Formulae 和 Casks，还能像原生系统服务一样管理终端工具的启动与停止。

---

## ✨ 核心特性

- 🖥️ **原生视觉体验**：采用 macOS Frameless 窗口设计，完美支持 **Vibrancy（毛玻璃）** 效果。
- ⚡ **服务一键管理**：支持对 `brew services` 进行图形化操作（启动/停止/重启）。
- 🔍 **丝滑搜索**：实时过滤海量软件清单，瞬间找到目标。
- ⏱️ **智能刷新**：每 2 分钟自动刷新当前 tab 数据，底部进度条显示剩余时间。
- 📦 **双列表展示**：清晰区分终端工具 (Formulae) 与桌面应用 (Casks)。
- 🔌 **Tap 管理**：全面支持 Homebrew Tap 的添加、移除、更新操作，轻松管理第三方软件源。
- 🐳 **Docker 容器管理**：查看、启动、停止 Docker 容器，实时监控容器状态。
- 📊 **侧边栏统计**：实时显示 GUI 应用、命令行工具、Taps、Docker 容器数量。

---

## 🆕 v1.4.0 新功能

### 🐳 Docker 容器管理
- 查看所有 Docker 容器列表（运行状态、镜像、端口）
- 一键启动/停止容器
- 本地搜索过滤容器
- 实时状态指示（运行中/已停止/暂停）
- 自动检测 Docker 是否安装/启动
- 侧边栏实时显示容器数量徽章
- 参与自动刷新（每 2 分钟更新）

---

## 🆕 v1.3.0 新功能

### 🔌 Tap 管理
- 完整的 Tap 管理界面
- 查看已安装的 Taps（官方/第三方分类）
- 添加新的 Tap
- 移除第三方 Tap（官方 Tap 保护机制）
- 批量更新所有 Taps
- Tap 详情查看（Git URL、描述）

### 🔍 软件包搜索
- 支持搜索 Formulae、Casks、Taps
- 一键安装搜索结果
- 从搜索结果直接添加第三方 Tap

### ⚡ 性能优化
- 🚀 **智能刷新**：只刷新当前 tab 数据，不刷新全量
- ⏱️ **进度条显示**：底部进度条实时显示下次刷新时间
- 📊 **侧边栏统计**：实时显示各类别数量徽章

### 🎨 UI 改进
- 侧边栏显示各类别数量统计徽章
- 专业 SVG 图标替换 Emoji
- 刷新按钮带旋转动画
- 毛玻璃效果增强

---

## 📸 界面预览

| 概览 | GUI 应用 | Tap 管理 |
| :--- | :--- | :--- |
| ![Overview](https://cdn.jsdelivr.net/gh/mycherish/imgCloud/img/20260318235836105.png) | ![GUI Apps](https://cdn.jsdelivr.net/gh/mycherish/imgCloud/img/20260318235927897.png) | ![Tap Management](https://cdn.jsdelivr.net/gh/mycherish/imgCloud/img/20260319000015839.png) |

---

## 🛠️ 技术栈

| 维度 | 技术 |
| :--- | :--- |
| **框架** | [Wails v2](https://wails.io/) (Go + WebView) |
| **前端** | Vue 3 (Composition API) + Vite |
| **样式** | CSS3 Glassmorphism + Native macOS Titlebar Inset |
| **后端** | Go (执行 brew 命令) |
| **平台** | macOS Apple Silicon (ARM64) |

---

## 🚀 快速开始

### 1. 前置要求
确保你的 Mac 已安装：
- [Go](https://go.dev/dl/) (1.20+)
- [Node.js](https://nodejs.org/) & [NPM](https://www.npmjs.com/)
- [Wails CLI](https://wails.io/docs/gettingstarted/installation)
- [Homebrew](https://brew.sh/)
- **Apple Silicon Mac** (M1/M2/M3/M4)

### 2. 开发模式
```bash
# 克隆仓库
git clone https://github.com/mycherish/brew-manager.git

# 进入目录
cd brew-manager

# 启动开发服务器
wails dev
```

---

## 💾 安装与使用

### 下载安装

1. 前往 [Releases](https://github.com/mycherish/brew-manager/releases) 页面下载最新的 `.dmg` 文件（仅支持 Apple Silicon）。
2. 打开 `.dmg` 并将 `Brew-Manager.app` 拖入 **Applications** 文件夹。

### ⚠️ 解决“无法验证开发者”问题
由于本应用未进行 Apple 开发者签名，首次打开时可能会提示“无法验证开发者”或“应用已损坏”。请执行以下操作：

1. **右键点击**应用程序文件夹中的 `Brew-Manager`，选择 **“打开”**，在弹出的对话框中再次点击 **“打开”**。
2. 如果依然无法运行，请打开终端执行以下命令：
   ```bash
   sudo xattr -rd com.apple.quarantine /Applications/Brew-Manager.app
   ```
3. 输入开机密码后即可正常使用。

---

## 📖 使用指南

### 1. 查看软件
- **概览**：查看所有软件统计
- **GUI 应用**：查看已安装的桌面应用
- **命令行工具**：查看已安装的终端工具

### 2. 服务管理
- 点击 **启动/停止** 按钮控制服务运行状态
- 点击 **重启** 按钮重启服务
- 状态指示灯实时显示服务状态

### 3. Tap 管理
- 点击 **+ 添加 Tap** 添加第三方软件源
- 点击 **🔄** 更新单个 Tap
- 点击 **更新所有** 批量更新所有 Taps
- 点击 **🗑️** 移除第三方 Tap（官方 Tap 有保护）

### 4. 软件包搜索
- 在 Taps 页面点击 **添加 Tap** 进入搜索模式
- 输入关键词搜索 Formulae、Casks、Taps
- 点击 **安装** 一键安装软件包

### 5. Docker 容器管理
- 查看所有 Docker 容器（运行中、已停止、暂停）
- 点击 **启动/停止** 按钮控制容器
- 使用搜索框过滤容器（按名称或镜像）
- 状态指示灯显示容器当前状态
- 点击刷新按钮更新容器列表

---

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

---

## 📄 许可证

MIT License

---

## 🙏 致谢

- [Wails](https://wails.io/) - 让 Go 桌面应用开发如此简单
- [Homebrew](https://brew.sh/) - macOS 包管理神器
- 所有开源贡献者

---

**如果你喜欢这个项目，欢迎 Star ⭐ 支持！**
