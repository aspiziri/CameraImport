package main

import (
	"context"
	"embed"
	"log"

	"cameraimport/internal/importer"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create application instance
	importService := importer.NewImportService()
	app := NewApp(importService)

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Camera Import",
		Width:  1400,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			importService.Startup(ctx)
		},
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
		},
		Bind: []interface{}{
			app,
			importService,
		},
	})

	if err != nil {
		log.Fatal("Error:", err.Error())
	}
}
