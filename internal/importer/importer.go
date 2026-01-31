package importer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"cameraimport/internal/config"
	"cameraimport/internal/exif"
	"cameraimport/internal/preview"
)

type FileInfo struct {
	Name      string            `json:"name"`
	Path      string            `json:"path"`
	Size      int64             `json:"size"`
	Date      string            `json:"date"`
	Type      string            `json:"type"`
	Thumbnail string            `json:"thumbnail,omitempty"`
	Exif      *exif.ExifData    `json:"exif,omitempty"`
	VideoMeta *VideoMetadata    `json:"videoMeta,omitempty"`
}

type VideoMetadata struct {
	Duration   float64 `json:"duration"`
	Width      int     `json:"width"`
	Height     int     `json:"height"`
	FrameRate  string  `json:"frameRate"`
	Resolution string  `json:"resolution"`
}

type ImportProgress struct {
	Current      int    `json:"current"`
	Total        int    `json:"total"`
	Percentage   int    `json:"percentage"`
	CurrentFile  string `json:"currentFile"`
	Status       string `json:"status"`
}

type ImportService struct {
	ctx       context.Context
	config    *config.Config
	progress  *ImportProgress
	mu        sync.RWMutex
	cancelled bool
	cancelMu  sync.RWMutex
}

func NewImportService() *ImportService {
	return &ImportService{
		config: config.NewConfig(),
		progress: &ImportProgress{
			Status:     "idle",
			Percentage: 0,
			Current:    0,
			Total:      0,
		},
	}
}

// Startup is called when the app starts
func (s *ImportService) Startup(ctx context.Context) {
	s.ctx = ctx
}

// ScanFiles scans the source directory for media files
// Fast scan - only gets basic file info, no thumbnails or EXIF during scan
func (s *ImportService) ScanFiles(sourcePath string, formats []string) ([]FileInfo, error) {
	var files []FileInfo

	fmt.Printf("Starting scan of: %s\n", sourcePath)

	err := filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %s: %v\n", path, err)
			return nil // Skip files with errors
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.TrimPrefix(filepath.Ext(path), ".")
		if !s.config.IsMediaFormat(ext) {
			return nil
		}

		// Use file modification time as date (faster than EXIF reading)
		dateStr := info.ModTime().Format("2006-01-02")

		fileInfo := FileInfo{
			Name: info.Name(),
			Path: path,
			Size: info.Size(),
			Date: dateStr,
		}

		// Determine file type (just by extension for speed)
		if s.config.IsImageFormat(ext) {
			fileInfo.Type = "image"
		} else if s.config.IsVideoFormat(ext) {
			fileInfo.Type = "video"
		} else if s.config.IsRawFormat(ext) {
			fileInfo.Type = "raw"
		}

		files = append(files, fileInfo)

		// Log progress every 100 files
		if len(files)%100 == 0 {
			fmt.Printf("Found %d files so far...\n", len(files))
		}

		return nil
	})

	fmt.Printf("Scan complete. Found %d files\n", len(files))
	return files, err
}

// GenerateThumbnail generates a thumbnail for a specific file path
func (s *ImportService) GenerateThumbnail(filePath string) (string, error) {
	thumbnail, err := preview.GenerateThumbnail(filePath)
	if err != nil {
		return "", err
	}
	return thumbnail, nil
}

// GenerateThumbnailsForFiles generates thumbnails for a batch of file paths
func (s *ImportService) GenerateThumbnailsForFiles(filePaths []string, maxCount int) map[string]string {
	thumbnails := make(map[string]string)
	count := 0

	for _, path := range filePaths {
		if count >= maxCount {
			break
		}

		// Check if it's an image file
		ext := strings.TrimPrefix(filepath.Ext(path), ".")
		if !s.config.IsImageFormat(ext) {
			continue
		}

		thumbnail := preview.GenerateThumbnailSafe(path)
		if thumbnail != "" {
			thumbnails[path] = thumbnail
			count++
		}

		if count%10 == 0 {
			fmt.Printf("Generated %d thumbnails...\n", count)
		}
	}

	return thumbnails
}

// RunImport performs the actual import operation asynchronously
func (s *ImportService) RunImport(sourcePath, destPath string, dates []string, folderName string, deleteAfter bool) error {
	// Reset cancelled flag
	s.cancelMu.Lock()
	s.cancelled = false
	s.cancelMu.Unlock()

	// Reset progress with mutex
	s.mu.Lock()
	s.progress.Status = "starting"
	s.progress.Current = 0
	s.progress.Total = 0
	s.progress.Percentage = 0
	s.progress.CurrentFile = ""
	s.mu.Unlock()

	fmt.Printf("Starting import: source=%s, dest=%s, dates=%v\n", sourcePath, destPath, dates)

	// Run import in goroutine so UI doesn't freeze
	go func() {
		files, err := s.ScanFiles(sourcePath, nil)
		if err != nil {
			s.mu.Lock()
			s.progress.Status = "error"
			s.mu.Unlock()
			fmt.Printf("Error scanning files: %v\n", err)
			return
		}

		// Filter files by selected dates if provided
		var filesToImport []FileInfo
		if len(dates) > 0 {
			dateMap := make(map[string]bool)
			for _, date := range dates {
				dateMap[date] = true
			}
			for _, file := range files {
				if dateMap[file.Date] {
					filesToImport = append(filesToImport, file)
				}
			}
		} else {
			filesToImport = files
		}

		fmt.Printf("Will import %d files\n", len(filesToImport))

		s.mu.Lock()
		s.progress.Total = len(filesToImport)
		s.progress.Current = 0
		s.progress.Status = "running"
		s.mu.Unlock()

		for i, file := range filesToImport {
			// Check for cancellation
			if s.isCancelled() {
				fmt.Println("Import cancelled, stopping...")
				s.mu.Lock()
				s.progress.Status = "cancelled"
				s.mu.Unlock()
				return
			}

			s.mu.Lock()
			s.progress.Current = i + 1
			if len(filesToImport) > 0 {
				s.progress.Percentage = int(float64(i+1) / float64(len(filesToImport)) * 100)
			}
			s.progress.CurrentFile = file.Name
			s.mu.Unlock()

			fmt.Printf("Importing %d/%d: %s\n", i+1, len(filesToImport), file.Name)

			if err := s.importFile(file, destPath, folderName); err != nil {
				fmt.Printf("Error importing %s: %v\n", file.Name, err)
				// Continue with other files
			}
		}

		if deleteAfter {
			fmt.Println("Deleting source files...")
			for _, file := range filesToImport {
				os.Remove(file.Path)
			}
		}

		s.mu.Lock()
		s.progress.Status = "completed"
		s.progress.Percentage = 100
		s.mu.Unlock()

		fmt.Println("Import completed!")
	}()

	return nil
}

func (s *ImportService) importFile(file FileInfo, destBase string, customFolder string) error {
	var destDir string

	if customFolder != "" {
		destDir = filepath.Join(destBase, customFolder)
	} else {
		destDir = filepath.Join(destBase, file.Date)
	}

	// Determine subdirectory based on file type
	var subDir string
	switch file.Type {
	case "image":
		subDir = strings.TrimPrefix(s.config.ImgRelativePath, "/")
	case "video":
		subDir = strings.TrimPrefix(s.config.VideoRelativePath, "/")
	case "raw":
		subDir = strings.TrimPrefix(s.config.RawRelativePath, "/")
	}

	fullDestDir := filepath.Join(destDir, subDir)
	if err := os.MkdirAll(fullDestDir, 0755); err != nil {
		return err
	}

	// Generate destination filename
	captureDate := exif.GetCaptureDate(file.Path)
	ext := filepath.Ext(file.Name)

	var destFileName string
	if file.Type == "video" && file.VideoMeta != nil {
		// Video filename with metadata
		destFileName = fmt.Sprintf("%s - %s - %s - %dsec%s",
			captureDate.Format("2006-01-02 15.04.05"),
			file.VideoMeta.Resolution,
			file.VideoMeta.FrameRate,
			int(file.VideoMeta.Duration),
			ext,
		)
	} else {
		// Standard filename
		destFileName = captureDate.Format("2006-01-02 15.04.05") + ext
	}

	destPath := filepath.Join(fullDestDir, destFileName)

	// Check if file already exists
	if _, err := os.Stat(destPath); err == nil {
		// File exists, check if it's the same
		if filesAreIdentical(file.Path, destPath) {
			return nil // Skip identical file
		}
		// Add microseconds to make filename unique
		destFileName = captureDate.Format("2006-01-02 15.04.05.000000") + ext
		destPath = filepath.Join(fullDestDir, destFileName)
	}

	// Copy the file
	return copyFile(file.Path, destPath)
}

func (s *ImportService) GetProgress() *ImportProgress {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return a copy to avoid race conditions
	progressCopy := *s.progress
	return &progressCopy
}

// CancelImport cancels the current import operation
func (s *ImportService) CancelImport() error {
	s.cancelMu.Lock()
	s.cancelled = true
	s.cancelMu.Unlock()

	s.mu.Lock()
	s.progress.Status = "cancelled"
	s.mu.Unlock()

	fmt.Println("Import cancelled by user")
	return nil
}

// isCancelled checks if the import has been cancelled
func (s *ImportService) isCancelled() bool {
	s.cancelMu.RLock()
	defer s.cancelMu.RUnlock()
	return s.cancelled
}

// Helper functions

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	// Preserve modification time
	if info, err := os.Stat(src); err == nil {
		os.Chtimes(dst, time.Now(), info.ModTime())
	}

	return nil
}

func filesAreIdentical(path1, path2 string) bool {
	file1, err := os.Open(path1)
	if err != nil {
		return false
	}
	defer file1.Close()

	file2, err := os.Open(path2)
	if err != nil {
		return false
	}
	defer file2.Close()

	info1, _ := file1.Stat()
	info2, _ := file2.Stat()

	if info1.Size() != info2.Size() {
		return false
	}

	// For performance, only compare file sizes
	// Could add hash comparison for more accuracy
	return true
}

func getVideoMetadata(path string) (*VideoMetadata, error) {
	// Use ffprobe to get video metadata
	cmd := exec.Command("ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		path,
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var result struct {
		Streams []struct {
			Width       int    `json:"width"`
			Height      int    `json:"height"`
			AvgFrameRate string `json:"avg_frame_rate"`
			Duration    string `json:"duration"`
		} `json:"streams"`
		Format struct {
			Duration string `json:"duration"`
		} `json:"format"`
	}

	if err := json.Unmarshal(output, &result); err != nil {
		return nil, err
	}

	if len(result.Streams) == 0 {
		return nil, fmt.Errorf("no video streams found")
	}

	stream := result.Streams[0]
	meta := &VideoMetadata{
		Width:      stream.Width,
		Height:     stream.Height,
		Resolution: fmt.Sprintf("%dw", stream.Width),
	}

	// Parse frame rate
	if parts := strings.Split(stream.AvgFrameRate, "/"); len(parts) == 2 {
		var num, denom float64
		fmt.Sscanf(parts[0], "%f", &num)
		fmt.Sscanf(parts[1], "%f", &denom)
		if denom > 0 {
			fps := num / denom
			meta.FrameRate = fmt.Sprintf("%.3f", fps)
		}
	}

	// Parse duration
	if stream.Duration != "" {
		fmt.Sscanf(stream.Duration, "%f", &meta.Duration)
	} else if result.Format.Duration != "" {
		fmt.Sscanf(result.Format.Duration, "%f", &meta.Duration)
	}

	return meta, nil
}
