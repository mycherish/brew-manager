# 🍺 Brew Manager

> **基于 Wails + Vue 3 构建的 macOS 原生质感 Homebrew 管理工具**

[![Wails](https://img.shields.io/badge/Built%20with-Wails-red.svg)](https://wails.io/)
[![Vue](https://img.shields.io/badge/Frontend-Vue%203-brightgreen.svg)](https://vuejs.org/)
[![Go](https://img.shields.io/badge/Backend-Go-blue.svg)](https://go.dev/)
[![Platform](https://img.shields.io/badge/Platform-macOS-lightgrey.svg)]()

`Brew Manager` 是一个为 macOS 用户打造的轻量级 Homebrew 图形界面工具。它不仅能让你一眼看清系统安装的所有 Formulae 和 Casks，还能像原生系统服务一样管理终端工具的启动与停止。

---

## ✨ 核心特性

- 🖥️ **原生视觉体验**：采用 macOS Frameless 窗口设计，完美支持 **Vibrancy（毛玻璃）** 效果。
- ⚡ **服务一键管理**：支持对 `brew services` 进行图形化操作（启动/停止）。
- 🔍 **丝滑搜索**：实时过滤海量软件清单，瞬间找到目标。
- 🔄 **自动同步**：每 10 秒自动更新软件运行状态，无需手动刷新。
- 📦 **双列表展示**：清晰区分终端工具 (Formulae) 与桌面应用 (Casks)。

---

## 📸 界面预览

| 概览 | 搜索与操作 |
| :--- | :--- |
| ![Main Window](https://cdn.jsdelivr.net/gh/mycherish/imgCloud/img/20260209115118699.png) | ![Search](https://cdn.jsdelivr.net/gh/mycherish/imgCloud/img/20260209111040397.png) |

---

## 🛠️ 技术栈

| 维度 | 技术 |
| :--- | :--- |
| **框架** | [Wails v2](https://wails.io/) (Go + Webview) |
| **前端** | Vue 3 (Composition API) + Vite |
| **样式** | CSS3 Glassmorphism + Native macOS Titlebar Inset |
| **后端** | Go (executing brew commands) |

---

## 🚀 快速开始

### 1. 前置要求
确保你的 Mac 已安装：
- [Go](https://go.dev/dl/) (1.20+)
- [Node.js](https://nodejs.org/) & [NPM](https://www.npmjs.com/)
- [Wails CLI](https://wails.io/docs/gettingstarted/installation)
- [Homebrew](https://brew.sh/)

### 2. 开发模式
```bash
# 克隆仓库
git clone https://github.com/mycherish/brew-manager.git

# 进入目录
cd brew-manager

# 启动开发服务器
wails dev