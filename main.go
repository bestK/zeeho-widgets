package main

import (
	"context"
	"embed"

	"github.com/bestk/zeeho-widgets/backend"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:  "Zeeho Widget",
		Width:  440,
		Height: 300,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		Windows: &windows.Options{
			WebviewIsTransparent:              true,
			WindowIsTranslucent:               false,
			DisableFramelessWindowDecorations: true,
		},
		OnStartup: func(ctx context.Context) {

			app.startup(ctx)

			backend.SetTransparentBackground()

			app.ScheduleRefresh()

		},
		Bind: []interface{}{
			app,
		},
		Frameless:         true,
		AlwaysOnTop:       false,
		DisableResize:     false,
		StartHidden:       false,
		HideWindowOnClose: false,
		WindowStartState:  options.Normal,
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
