package exif

import (
	"fmt"
	"os"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

type ExifData struct {
	DateTime     time.Time `json:"dateTime"`
	Camera       string    `json:"camera,omitempty"`
	Lens         string    `json:"lens,omitempty"`
	ISO          int       `json:"iso,omitempty"`
	FNumber      string    `json:"fNumber,omitempty"`
	ExposureTime string    `json:"exposureTime,omitempty"`
	FocalLength  string    `json:"focalLength,omitempty"`
	Width        int       `json:"width,omitempty"`
	Height       int       `json:"height,omitempty"`
	Dimensions   string    `json:"dimensions,omitempty"`
}

func ReadExif(filePath string) (*ExifData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	x, err := exif.Decode(file)
	if err != nil {
		return nil, err
	}

	data := &ExifData{}

	// Extract date/time
	if dt, err := x.DateTime(); err == nil {
		data.DateTime = dt
	}

	// Extract camera make and model
	if make, err := x.Get(exif.Make); err == nil {
		if makeStr, err := make.StringVal(); err == nil {
			if model, err := x.Get(exif.Model); err == nil {
				if modelStr, err := model.StringVal(); err == nil {
					data.Camera = fmt.Sprintf("%s %s", makeStr, modelStr)
				}
			}
		}
	}

	// Extract lens info
	if lens, err := x.Get(exif.LensModel); err == nil {
		if lensStr, err := lens.StringVal(); err == nil {
			data.Lens = lensStr
		}
	}

	// Extract ISO
	if iso, err := x.Get(exif.ISOSpeedRatings); err == nil {
		if val, err := iso.Int(0); err == nil {
			data.ISO = val
		}
	}

	// Extract F-Number
	if fnum, err := x.Get(exif.FNumber); err == nil {
		if rat, err := fnum.Rat(0); err == nil {
			num := float64(rat.Num().Int64())
			denom := float64(rat.Denom().Int64())
			if denom > 0 {
				data.FNumber = fmt.Sprintf("f/%.1f", num/denom)
			}
		}
	}

	// Extract exposure time
	if expTime, err := x.Get(exif.ExposureTime); err == nil {
		if rat, err := expTime.Rat(0); err == nil {
			data.ExposureTime = fmt.Sprintf("%d/%d", rat.Num().Int64(), rat.Denom().Int64())
		}
	}

	// Extract focal length
	if focal, err := x.Get(exif.FocalLength); err == nil {
		if rat, err := focal.Rat(0); err == nil {
			num := float64(rat.Num().Int64())
			denom := float64(rat.Denom().Int64())
			if denom > 0 {
				data.FocalLength = fmt.Sprintf("%.0fmm", num/denom)
			}
		}
	}

	// Extract image dimensions
	if width, err := x.Get(exif.PixelXDimension); err == nil {
		if w, err := width.Int(0); err == nil {
			data.Width = w
		}
	}
	if height, err := x.Get(exif.PixelYDimension); err == nil {
		if h, err := height.Int(0); err == nil {
			data.Height = h
		}
	}

	if data.Width > 0 && data.Height > 0 {
		data.Dimensions = fmt.Sprintf("%dx%d", data.Width, data.Height)
	}

	return data, nil
}

// GetCaptureDate tries to get the actual capture date from EXIF,
// falls back to file modification time if EXIF is not available
func GetCaptureDate(filePath string) time.Time {
	if exifData, err := ReadExif(filePath); err == nil && !exifData.DateTime.IsZero() {
		return exifData.DateTime
	}

	return FileModDate(filePath)
}

// FileModDate returns the file's modification time, or the current time if
// the file cannot be read
func FileModDate(filePath string) time.Time {
	if info, err := os.Stat(filePath); err == nil {
		return info.ModTime()
	}

	return time.Now()
}
