//go:build darwin

package usbdetector

/*
#cgo LDFLAGS: -framework DiskArbitration -framework CoreFoundation
#include <DiskArbitration/DiskArbitration.h>
#include <stdlib.h>

extern void diskAppearedCallback(DADiskRef disk, void *context);

static inline void startSession(DASessionRef session, CFRunLoopRef runLoop) {
    DASessionScheduleWithRunLoop(session, runLoop, kCFRunLoopDefaultMode);
}

static inline void registerCallback(DASessionRef session, void *context) {
    DARegisterDiskAppearedCallback(session, NULL, diskAppearedCallback, context);
}

static inline char* getDiskPath(DADiskRef disk) {
    CFDictionaryRef desc = DADiskCopyDescription(disk);
    if (desc == NULL) return NULL;

    CFURLRef pathURL = CFDictionaryGetValue(desc, kDADiskDescriptionVolumePathKey);
    if (pathURL == NULL) {
        CFRelease(desc);
        return NULL;
    }

    CFStringRef pathString = CFURLCopyFileSystemPath(pathURL, kCFURLPOSIXPathStyle);
    if (pathString == NULL) {
        CFRelease(desc);
        return NULL;
    }

    CFIndex length = CFStringGetLength(pathString);
    CFIndex maxSize = CFStringGetMaximumSizeForEncoding(length, kCFStringEncodingUTF8) + 1;
    char *buffer = malloc(maxSize);

    if (CFStringGetCString(pathString, buffer, maxSize, kCFStringEncodingUTF8)) {
        CFRelease(pathString);
        CFRelease(desc);
        return buffer;
    }

    free(buffer);
    CFRelease(pathString);
    CFRelease(desc);
    return NULL;
}

static inline char* getDiskName(DADiskRef disk) {
    CFDictionaryRef desc = DADiskCopyDescription(disk);
    if (desc == NULL) return NULL;

    CFStringRef name = CFDictionaryGetValue(desc, kDADiskDescriptionVolumeNameKey);
    if (name == NULL) {
        CFRelease(desc);
        return NULL;
    }

    CFIndex length = CFStringGetLength(name);
    CFIndex maxSize = CFStringGetMaximumSizeForEncoding(length, kCFStringEncodingUTF8) + 1;
    char *buffer = malloc(maxSize);

    if (CFStringGetCString(name, buffer, maxSize, kCFStringEncodingUTF8)) {
        CFRelease(desc);
        return buffer;
    }

    free(buffer);
    CFRelease(desc);
    return NULL;
}

static inline uint64_t getDiskSize(DADiskRef disk) {
    CFDictionaryRef desc = DADiskCopyDescription(disk);
    if (desc == NULL) return 0;

    CFNumberRef sizeNum = CFDictionaryGetValue(desc, kDADiskDescriptionMediaSizeKey);
    if (sizeNum == NULL) {
        CFRelease(desc);
        return 0;
    }

    uint64_t size;
    CFNumberGetValue(sizeNum, kCFNumberSInt64Type, &size);
    CFRelease(desc);
    return size;
}

static inline int isRemovable(DADiskRef disk) {
    CFDictionaryRef desc = DADiskCopyDescription(disk);
    if (desc == NULL) return 0;

    CFBooleanRef removable = CFDictionaryGetValue(desc, kDADiskDescriptionMediaRemovableKey);
    int result = (removable != NULL && CFBooleanGetValue(removable));
    CFRelease(desc);
    return result;
}
*/
import "C"
import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"unsafe"
)

type DarwinDetector struct {
	ctx      context.Context
	cancel   context.CancelFunc
	callback func(DriveInfo)
	session  C.DASessionRef
	mu       sync.Mutex
}

var globalDetector *DarwinDetector

// NewDetector creates a new macOS USB detector
func NewDetector() Detector {
	return &DarwinDetector{}
}

// Start begins monitoring for USB device connections
func (d *DarwinDetector) Start(ctx context.Context, callback func(DriveInfo)) error {
	d.ctx, d.cancel = context.WithCancel(ctx)
	d.callback = callback
	globalDetector = d

	// Create a DiskArbitration session
	d.session = C.DASessionCreate(C.kCFAllocatorDefault)
	if d.session == nil {
		return fmt.Errorf("failed to create DiskArbitration session")
	}

	// Schedule with run loop
	runLoop := C.CFRunLoopGetCurrent()
	C.startSession(d.session, runLoop)

	// Register callback
	C.registerCallback(d.session, unsafe.Pointer(d))

	// Start run loop in goroutine
	go func() {
		C.CFRunLoopRun()
	}()

	// Stop when context is cancelled
	go func() {
		<-d.ctx.Done()
		C.CFRunLoopStop(C.CFRunLoopGetCurrent())
	}()

	return nil
}

// Stop stops monitoring for USB device connections
func (d *DarwinDetector) Stop() error {
	if d.cancel != nil {
		d.cancel()
	}

	if d.session != nil {
		C.DASessionUnscheduleFromRunLoop(d.session, C.CFRunLoopGetCurrent(), C.kCFRunLoopDefaultMode)
		C.CFRelease(C.CFTypeRef(d.session))
		d.session = nil
	}

	return nil
}

// GetRemovableDrives returns a list of currently connected removable drives
func (d *DarwinDetector) GetRemovableDrives() ([]DriveInfo, error) {
	drives := []DriveInfo{}

	volumesPath := "/Volumes"
	entries, err := os.ReadDir(volumesPath)
	if err != nil {
		return drives, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		path := filepath.Join(volumesPath, entry.Name())

		// Get filesystem info to check if removable
		var stat syscall.Statfs_t
		err := syscall.Statfs(path, &stat)
		if err != nil {
			continue
		}

		// Skip network and system volumes
		if stat.Flags&syscall.MNT_LOCAL == 0 {
			continue
		}

		info := DriveInfo{
			Name:        entry.Name(),
			Path:        path,
			Size:        uint64(stat.Blocks) * uint64(stat.Bsize),
			IsRemovable: true, // Volumes in /Volumes are typically removable
		}
		info.SizeGB = FormatSize(info.Size)

		drives = append(drives, info)
	}

	return drives, nil
}

// This is called from C when a disk appears
//export diskAppearedCallback
func diskAppearedCallback(disk C.DADiskRef, context unsafe.Pointer) {
	if globalDetector == nil || globalDetector.callback == nil {
		return
	}

	// Check if removable
	if C.isRemovable(disk) == 0 {
		return
	}

	// Get disk info
	pathC := C.getDiskPath(disk)
	if pathC == nil {
		return
	}
	defer C.free(unsafe.Pointer(pathC))

	nameC := C.getDiskName(disk)
	name := "Untitled"
	if nameC != nil {
		name = C.GoString(nameC)
		C.free(unsafe.Pointer(nameC))
	}

	path := C.GoString(pathC)
	size := uint64(C.getDiskSize(disk))

	driveInfo := DriveInfo{
		Name:        name,
		Path:        path,
		Size:        size,
		SizeGB:      FormatSize(size),
		IsRemovable: true,
	}

	globalDetector.callback(driveInfo)
}
