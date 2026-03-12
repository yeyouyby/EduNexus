package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the backend structure
	backend := NewBackend()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "EduNexus - 校园超算可视化中枢",
		Width:  1280,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 9, G: 13, B: 20, A: 255}, // #090D14
		Frameless:        true,
		OnStartup:        backend.startup,
		Bind: []interface{}{
			backend,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
