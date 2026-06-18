//go:build linux

package usbdetector

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

type LinuxDetector struct {
	ctx       context.Context
	cancel    context.CancelFunc
	callback  func(DriveInfo)
	lastCheck map[string]bool
}

// NewDetector creates a new Linux USB detector
func NewDetector() Detector {
	return &LinuxDetector{
		lastCheck: make(map[string]bool),
	}
}

// Start begins monitoring for USB device connections
func (d *LinuxDetector) Start(ctx context.Context, callback func(DriveInfo)) error {
	d.ctx, d.cancel = context.WithCancel(ctx)
	d.callback = callback

	// Initialize with current drives
	drives, err := d.GetRemovableDrives()
	if err == nil {
		for _, drive := range drives {
			d.lastCheck[drive.Path] = true
		}
	}

	// Start polling for changes
	go d.pollForChanges()

	return nil
}

// Stop stops monitoring for USB device connections
func (d *LinuxDetector) Stop() error {
	if d.cancel != nil {
		d.cancel()
	}
	return nil
}

// pollForChanges checks for new drives periodically
func (d *LinuxDetector) pollForChanges() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-d.ctx.Done():
			return
		case <-ticker.C:
			drives, err := d.GetRemovableDrives()
			if err != nil {
				continue
			}

			currentDrives := make(map[string]bool)
			for _, drive := range drives {
				currentDrives[drive.Path] = true

				// Check if this is a new drive
				if !d.lastCheck[drive.Path] {
					if d.callback != nil {
						d.callback(drive)
					}
				}
			}

			d.lastCheck = currentDrives
		}
	}
}

// GetRemovableDrives returns a list of currently connected removable drives
func (d *LinuxDetector) GetRemovableDrives() ([]DriveInfo, error) {
	drives := []DriveInfo{}

	// Check common mount points
	mountPoints := []string{"/media", "/mnt", "/run/media"}

	for _, mountBase := range mountPoints {
		if _, err := os.Stat(mountBase); os.IsNotExist(err) {
			continue
		}

		// For /run/media, check subdirectories (user-specific mounts)
		if mountBase == "/run/media" {
			userDirs, err := os.ReadDir(mountBase)
			if err != nil {
				continue
			}

			for _, userDir := range userDirs {
				if !userDir.IsDir() {
					continue
				}
				userPath := filepath.Join(mountBase, userDir.Name())
				d.scanMountPoint(userPath, &drives)
			}
		} else {
			d.scanMountPoint(mountBase, &drives)
		}
	}

	return drives, nil
}

// scanMountPoint scans a mount point for removable drives
func (d *LinuxDetector) scanMountPoint(mountBase string, drives *[]DriveInfo) {
	entries, err := os.ReadDir(mountBase)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		path := filepath.Join(mountBase, entry.Name())

		// Check if it's a valid mount point
		var stat syscall.Statfs_t
		if err := syscall.Statfs(path, &stat); err != nil {
			continue
		}

		// Skip if total size is 0
		totalSize := uint64(stat.Blocks) * uint64(stat.Bsize)
		if totalSize == 0 {
			continue
		}

		// Check if it's removable by checking /sys/block
		isRemovable := d.isRemovable(path)
		if !isRemovable {
			continue
		}

		name := entry.Name()
		// Try to get a better name from .disk_label if it exists
		labelPath := filepath.Join(path, ".disk_label")
		if labelData, err := os.ReadFile(labelPath); err == nil {
			if label := strings.TrimSpace(string(labelData)); label != "" {
				name = label
			}
		}

		info := DriveInfo{
			Name:        name,
			Path:        path,
			Size:        totalSize,
			SizeGB:      FormatSize(totalSize),
			IsRemovable: true,
		}

		*drives = append(*drives, info)
	}
}

// isRemovable checks if a mount point corresponds to a removable device
func (d *LinuxDetector) isRemovable(mountPath string) bool {
	// Read /proc/mounts to find the device
	mountsData, err := os.ReadFile("/proc/mounts")
	if err != nil {
		return false
	}

	lines := strings.Split(string(mountsData), "\n")
	var deviceName string

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[1] == mountPath {
			deviceName = fields[0]
			break
		}
	}

	if deviceName == "" {
		return false
	}

	deviceName = filepath.Base(deviceName)

	// Resolve the parent disk via sysfs: /sys/class/block/<partition> points
	// into the disk's directory (e.g. .../block/mmcblk0/mmcblk0p1), while a
	// whole disk resolves to a path whose parent is "block". Stripping
	// trailing digits would mangle names like mmcblk0p1 or nvme0n1p1.
	if resolved, err := filepath.EvalSymlinks(filepath.Join("/sys/class/block", deviceName)); err == nil {
		if parent := filepath.Base(filepath.Dir(resolved)); parent != "block" {
			deviceName = parent
		}
	}

	// Check /sys/block/<device>/removable
	removablePath := fmt.Sprintf("/sys/block/%s/removable", deviceName)
	data, err := os.ReadFile(removablePath)
	if err != nil {
		return false
	}

	return strings.TrimSpace(string(data)) == "1"
}
