package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac" // 引入 Mac 专用选项
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "wails-test",
		Width:  1024,
		Height: 768,
		// Frameless: true, // 1. 隐藏标准标题栏

		Mac: &mac.Options{
			TitleBar:             mac.TitleBarHiddenInset(), // 4. 隐藏标题文字，但保留左上角红绿灯并内嵌
			WindowIsTranslucent:  true,
			WebviewIsTransparent: true,
			Appearance:           mac.NSAppearanceNameDarkAqua, // 强制深色模式 (可选)
			About: &mac.AboutInfo{
				Title:   "wails-test",
				Message: "你的专属 Brew 管理工具",
			},
		},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		// 关键点 1：最外层使用 BackgroundColour 来控制 Webview 基础透明度
		// 设置 Alpha 为 0 表示完全透明
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
