package main

import (
	"embed"
	"time"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"

	"github.com/notebox/nbfm/pkg/config"
	"github.com/notebox/nbfm/pkg/file"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp(&config.Config{SyncInterval: 5 * time.Second})

	// Create application with options
	err := wails.Run(&options.App{
		Title:     "nbfm",
		Width:     1024,
		Height:    768,
		MinWidth:  300,
		MinHeight: 300,
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: file.NewFileServer(),
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.Startup,
		OnShutdown:       app.Shutdown,
		Bind: []interface{}{
			app,
		},
		Mac: &mac.Options{
			Preferences: &mac.Preferences{
				FullscreenEnabled: mac.Enabled,
			},
		},
		Debug: options.Debug{
			OpenInspectorOnStartup: false,
		},
		Logger:             app.Logger,
		LogLevelProduction: logger.INFO,
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
