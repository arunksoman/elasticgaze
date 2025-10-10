package main

import (
	"context"
	"embed"
	"os"
	"path/filepath"

	"elasticgaze/internal/logging"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/build
var assets embed.FS

func main() {
	// Initialize logger early for Wails
	appDataDir, err := os.UserConfigDir()
	if err != nil {
		panic("Failed to get user config directory: " + err.Error())
	}
	elasticGazeDir := filepath.Join(appDataDir, "elasticgaze")

	// Initialize Wails-compatible logger
	wailsLogger, err := logging.InitLogger(elasticGazeDir)
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err = wails.Run(&options.App{
		Title:            "ElasticGaze",
		Width:            800,
		Height:           600,
		Frameless:        true,
		WindowStartState: options.Normal,
		MinWidth:         800,
		MinHeight:        600,
		Logger:           wailsLogger, // Register our custom logger with Wails
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 1},
		OnStartup:        app.startup,
		OnShutdown: func(ctx context.Context) {
			// Ensure database connection is properly closed on shutdown
			if err := app.Close(); err != nil {
				logging.Error("Error closing database:", err.Error())
			}
		},
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		logging.Error("Error:", err.Error())
	}
}
