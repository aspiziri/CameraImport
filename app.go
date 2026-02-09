package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"cameraimport/internal/usbdetector"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Config represents the application configuration
type Config struct {
	Source            string   `json:"source"`
	Destination       string   `json:"destination"`
	ImgFormats        []string `json:"imgFormats"`
	VideoFormats      []string `json:"videoFormats"`
	RawFormats        []string `json:"rawFormats"`
	ImgRelativePath   string   `json:"imgRelativePath"`
	VideoRelativePath string   `json:"videoRelativePath"`
	RawRelativePath   string   `json:"rawRelativePath"`
}

// App struct
type App struct {
	ctx         context.Context
	usbDetector usbdetector.Detector
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		usbDetector: usbdetector.NewDetector(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Start USB detection
	err := a.usbDetector.Start(ctx, func(drive usbdetector.DriveInfo) {
		log.Printf("USB drive detected: %s (%s) at %s", drive.Name, drive.SizeGB, drive.Path)

		// Emit event to frontend
		runtime.EventsEmit(ctx, "usb-device-connected", drive)
	})

	if err != nil {
		log.Printf("Failed to start USB detector: %v", err)
	}
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, CameraImport is ready!", name)
}

// SelectFolder opens a directory selection dialog
func (a *App) SelectFolder(title string) (string, error) {
	options := runtime.OpenDialogOptions{
		Title: title,
	}

	selected, err := runtime.OpenDirectoryDialog(a.ctx, options)
	if err != nil {
		return "", err
	}

	return selected, nil
}

// getConfigPath returns the path to the config file
func (a *App) getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".cameraimport")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(configDir, "config.json"), nil
}

// LoadConfig loads the saved configuration
func (a *App) LoadConfig() (*Config, error) {
	configPath, err := a.getConfigPath()
	if err != nil {
		return a.getDefaultConfig(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		// If file doesn't exist, return default config
		if os.IsNotExist(err) {
			return a.getDefaultConfig(), nil
		}
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return a.getDefaultConfig(), nil
	}

	return &config, nil
}

// SaveConfig saves the configuration to disk
func (a *App) SaveConfig(config Config) error {
	configPath, err := a.getConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// getDefaultConfig returns the default configuration
func (a *App) getDefaultConfig() *Config {
	return &Config{
		Source:            "",
		Destination:       "",
		ImgFormats:        []string{"jpg", "JPG", "png", "PNG", "gif", "GIF", "tif", "TIF", "jpeg", "JPEG"},
		VideoFormats:      []string{"mp4", "MP4", "avi", "AVI", "mov", "MOV"},
		RawFormats:        []string{"arw", "ARW", "cr2", "CR2", "nef", "NEF", "dng", "DNG"},
		ImgRelativePath:   "/",
		VideoRelativePath: "/Videos/",
		RawRelativePath:   "/Capture/",
	}
}

// GetRemovableDrives returns a list of currently connected removable drives
func (a *App) GetRemovableDrives() ([]usbdetector.DriveInfo, error) {
	return a.usbDetector.GetRemovableDrives()
}
