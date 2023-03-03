package main

import (
	"embed"
	"runtime/debug"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

func main() {
	// 设置GC限制为300M，触发到阈值后开始一轮GC
	// 因为文件传输经常会遇到大文件内存占用，而长时间不触发GC，可能会让机器内存榨干
	debug.SetMemoryLimit(1024 * 1024 * 300)
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:             "Ding",
		Width:             350,
		Height:            680,
		DisableResize:     true,
		Fullscreen:        false,
		StartHidden:       true,
		HideWindowOnClose: false,
		Frameless:         true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: options.NewRGBA(255, 255, 255, 0), // 先让背景色透明，然后在页面加载后显示窗口
		OnStartup:        app.startup,
		OnDomReady:       app.domReady,
		OnShutdown:       app.shutdown,
		Bind: []interface{}{
			app,
		},
		Mac: &mac.Options{
			TitleBar:             mac.TitleBarHiddenInset(),
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			About: &mac.AboutInfo{
				Title:   "Ding",
				Message: "© LionOcean 2023 ",
				Icon:    icon,
			},
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
