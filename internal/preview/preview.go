package preview

import (
	"encoding/base64"
	"bytes"
	"image"
	"image/jpeg"
	"os"

	"github.com/disintegration/imaging"
)

const (
	ThumbnailWidth  = 400
	ThumbnailHeight = 300
)

// GenerateThumbnail creates a thumbnail for an image file
func GenerateThumbnail(filePath string) (string, error) {
	// Open the image file
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	// Create thumbnail
	thumb := imaging.Fit(img, ThumbnailWidth, ThumbnailHeight, imaging.Lanczos)

	// Encode to JPEG in memory
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, thumb, &jpeg.Options{Quality: 80}); err != nil {
		return "", err
	}

	// Convert to base64 data URL
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	return "data:image/jpeg;base64," + encoded, nil
}

// GenerateThumbnailSafe attempts to generate a thumbnail, returns empty string on error
func GenerateThumbnailSafe(filePath string) string {
	thumbnail, err := GenerateThumbnail(filePath)
	if err != nil {
		return ""
	}
	return thumbnail
}
