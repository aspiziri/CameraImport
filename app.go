package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"cameraimport/internal/config"
	"cameraimport/internal/importer"
	"cameraimport/internal/usbdetector"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx           context.Context
	importService *importer.ImportService
}

// NewApp creates a new App application struct
func NewApp(importService *importer.ImportService) *App {
	return &App{importService: importService}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Apply the persisted configuration to the import service
	a.LoadConfig()
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

// LoadConfig loads the saved configuration and applies it to the import service
func (a *App) LoadConfig() (*config.Config, error) {
	cfg := a.loadConfigOrDefault()
	a.importService.SetConfig(cfg)
	return cfg, nil
}

func (a *App) loadConfigOrDefault() *config.Config {
	configPath, err := a.getConfigPath()
	if err != nil {
		return config.NewConfig()
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		// Missing or unreadable config falls back to defaults
		return config.NewConfig()
	}

	return cfg
}

// SaveConfig saves the configuration to disk and applies it to the import service
func (a *App) SaveConfig(cfg config.Config) error {
	configPath, err := a.getConfigPath()
	if err != nil {
		return err
	}

	if err := cfg.SaveConfig(configPath); err != nil {
		return err
	}

	a.importService.SetConfig(&cfg)
	return nil
}

// GetRemovableDrives returns a list of currently connected removable drives
func (a *App) GetRemovableDrives() ([]usbdetector.DriveInfo, error) {
	return usbdetector.NewDetector().GetRemovableDrives()
}
