//go:build windows

package usbdetector

import (
	"context"
	"fmt"
	"strconv"
	"syscall"
	"time"
	"unsafe"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

type WindowsDetector struct {
	ctx       context.Context
	cancel    context.CancelFunc
	callback  func(DriveInfo)
	lastCheck map[string]bool
}

// NewDetector creates a new Windows USB detector
func NewDetector() Detector {
	return &WindowsDetector{
		lastCheck: make(map[string]bool),
	}
}

// Start begins monitoring for USB device connections
func (d *WindowsDetector) Start(ctx context.Context, callback func(DriveInfo)) error {
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
func (d *WindowsDetector) Stop() error {
	if d.cancel != nil {
		d.cancel()
	}
	return nil
}

// pollForChanges checks for new drives periodically
func (d *WindowsDetector) pollForChanges() {
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
func (d *WindowsDetector) GetRemovableDrives() ([]DriveInfo, error) {
	drives := []DriveInfo{}

	// Initialize COM
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	// Connect to WMI
	unknown, err := oleutil.CreateObject("WbemScripting.SWbemLocator")
	if err != nil {
		return drives, err
	}
	defer unknown.Release()

	wmi, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return drives, err
	}
	defer wmi.Release()

	// Connect to WMI service
	serviceRaw, err := oleutil.CallMethod(wmi, "ConnectServer")
	if err != nil {
		return drives, err
	}
	service := serviceRaw.ToIDispatch()
	defer service.Release()

	// Query for logical disks with media present (Size IS NOT NULL excludes empty card reader slots)
	resultRaw, err := oleutil.CallMethod(service, "ExecQuery",
		"SELECT * FROM Win32_LogicalDisk WHERE DriveType = 2 AND Size IS NOT NULL")
	if err != nil {
		return drives, err
	}
	result := resultRaw.ToIDispatch()
	defer result.Release()

	// Get count of items
	countVar, err := oleutil.GetProperty(result, "Count")
	if err != nil {
		return drives, err
	}
	count := int(countVar.Val)

	// Iterate through results
	for i := 0; i < count; i++ {
		itemRaw, err := oleutil.CallMethod(result, "ItemIndex", i)
		if err != nil {
			continue
		}
		item := itemRaw.ToIDispatch()

		// Get drive properties
		deviceID, _ := oleutil.GetProperty(item, "DeviceID")
		volumeName, _ := oleutil.GetProperty(item, "VolumeName")
		size, _ := oleutil.GetProperty(item, "Size")

		driveInfo := DriveInfo{
			Path:        deviceID.ToString(),
			Name:        volumeName.ToString(),
			IsRemovable: true,
		}

		// WMI returns uint64 Size as a BSTR string, not a numeric VARIANT —
		// size.Val would be the string pointer address, not the byte count.
		if sizeStr := size.ToString(); sizeStr != "" {
			if sizeBytes, err := strconv.ParseUint(sizeStr, 10, 64); err == nil {
				driveInfo.Size = sizeBytes
				driveInfo.SizeGB = FormatSize(sizeBytes)
			}
		}

		if driveInfo.Name == "" {
			driveInfo.Name = fmt.Sprintf("Removable Disk (%s)", driveInfo.Path)
		}

		item.Release()
		drives = append(drives, driveInfo)
	}

	return drives, nil
}

// GetDriveType is a helper to check drive type using Windows API
func getDriveType(rootPath string) uint32 {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getDriveTypeW := kernel32.NewProc("GetDriveTypeW")

	rootPathPtr, _ := syscall.UTF16PtrFromString(rootPath)
	ret, _, _ := getDriveTypeW.Call(uintptr(unsafe.Pointer(rootPathPtr)))

	return uint32(ret)
}
