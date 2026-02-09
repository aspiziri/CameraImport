package usbdetector

import (
	"context"
)

// DriveInfo contains information about a detected drive
type DriveInfo struct {
	Name        string `json:"name"`        // Drive name/label
	Path        string `json:"path"`        // Mount point or drive letter
	Size        uint64 `json:"size"`        // Size in bytes
	SizeGB      string `json:"sizeGB"`      // Human-readable size
	IsRemovable bool   `json:"isRemovable"` // Whether it's a removable drive
}

// Detector is the interface for platform-specific USB detection
type Detector interface {
	// Start begins monitoring for USB device connections
	Start(ctx context.Context, callback func(DriveInfo)) error

	// Stop stops monitoring for USB device connections
	Stop() error

	// GetRemovableDrives returns a list of currently connected removable drives
	GetRemovableDrives() ([]DriveInfo, error)
}

// FormatSize converts bytes to human-readable format
func FormatSize(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return "< 1 GB"
	}

	gb := float64(bytes) / float64(unit*unit*unit)
	return formatFloat(gb) + " GB"
}

func formatFloat(f float64) string {
	if f < 10 {
		return formatTwoDecimals(f)
	}
	return formatOneDecimal(f)
}

func formatTwoDecimals(f float64) string {
	result := make([]byte, 0, 8)
	intPart := int(f)
	result = appendInt(result, intPart)
	result = append(result, '.')

	decimal := int((f - float64(intPart)) * 100)
	if decimal < 10 {
		result = append(result, '0')
	}
	result = appendInt(result, decimal)
	return string(result)
}

func formatOneDecimal(f float64) string {
	result := make([]byte, 0, 8)
	intPart := int(f)
	result = appendInt(result, intPart)
	result = append(result, '.')

	decimal := int((f - float64(intPart)) * 10)
	result = appendInt(result, decimal)
	return string(result)
}

func appendInt(b []byte, i int) []byte {
	if i == 0 {
		return append(b, '0')
	}

	if i < 0 {
		b = append(b, '-')
		i = -i
	}

	var tmp [20]byte
	j := len(tmp)
	for i > 0 {
		j--
		tmp[j] = byte('0' + i%10)
		i /= 10
	}

	return append(b, tmp[j:]...)
}
