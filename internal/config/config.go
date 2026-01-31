package config

import (
	"encoding/json"
	"os"
	"strings"
)

type Config struct {
	ImgFormats        []string `json:"imgFormats"`
	VideoFormats      []string `json:"videoFormats"`
	RawFormats        []string `json:"rawFormats"`
	Source            string   `json:"source"`
	Destination       string   `json:"destination"`
	ImgRelativePath   string   `json:"imgRelativePath"`
	VideoRelativePath string   `json:"videoRelativePath"`
	RawRelativePath   string   `json:"rawRelativePath"`
}

func NewConfig() *Config {
	return &Config{
		ImgFormats:        []string{"jpg", "JPG", "png", "PNG", "gif", "GIF", "tif", "TIF", "jpeg", "JPEG"},
		VideoFormats:      []string{"mp4", "MP4", "avi", "AVI", "mov", "MOV"},
		RawFormats:        []string{"arw", "ARW", "cr2", "CR2", "nef", "NEF", "dng", "DNG"},
		ImgRelativePath:   "/",
		VideoRelativePath: "/Videos/",
		RawRelativePath:   "/Capture/",
	}
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) SaveConfig(path string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func (c *Config) IsImageFormat(ext string) bool {
	return c.containsFormat(ext, c.ImgFormats)
}

func (c *Config) IsVideoFormat(ext string) bool {
	return c.containsFormat(ext, c.VideoFormats)
}

func (c *Config) IsRawFormat(ext string) bool {
	return c.containsFormat(ext, c.RawFormats)
}

func (c *Config) IsMediaFormat(ext string) bool {
	return c.IsImageFormat(ext) || c.IsVideoFormat(ext) || c.IsRawFormat(ext)
}

func (c *Config) containsFormat(ext string, formats []string) bool {
	ext = strings.TrimPrefix(ext, ".")
	for _, format := range formats {
		if strings.EqualFold(ext, format) {
			return true
		}
	}
	return false
}
