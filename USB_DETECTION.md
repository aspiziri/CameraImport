# USB Device Detection

## Overview

The Camera Import application now includes automatic USB device detection that alerts you when a removable drive (like a camera's SD card or USB flash drive) is connected to your computer.

## How It Works

### Platform-Specific Implementation

The USB detection system uses different approaches for each operating system:

**Windows:**
- Uses WMI (Windows Management Instrumentation) to detect removable drives
- Monitors for `Win32_LogicalDisk` with `DriveType = 2` (removable)
- Polls for changes every 2 seconds

**macOS:**
- Uses DiskArbitration framework (native macOS API)
- Real-time event-driven detection
- Monitors `/Volumes` directory for removable drives

**Linux:**
- Monitors common mount points (`/media`, `/mnt`, `/run/media`)
- Checks `/sys/block/<device>/removable` to verify if drive is removable
- Polls for changes every 2 seconds

### User Experience

1. **Automatic Detection:** When you plug in a USB device while the app is running, it's automatically detected
2. **Modal Prompt:** A modal dialog appears showing:
   - Device name
   - Mount point/drive letter
   - Total size in GB
3. **User Choice:** You can choose to:
   - **Yes, Use This Device:** Sets the device as your import source location
   - **No, Thanks:** Dismisses the dialog without changes

## Code Structure

```
internal/usbdetector/
├── detector.go           # Common interface and types
├── detector_windows.go   # Windows WMI implementation
├── detector_darwin.go    # macOS DiskArbitration implementation
└── detector_linux.go     # Linux udev/mount monitoring implementation
```

## API

### Backend (Go)

**Starting Detection:**
```go
detector := usbdetector.NewDetector()
err := detector.Start(ctx, func(drive usbdetector.DriveInfo) {
    // Called when a new USB device is detected
})
```

**Getting Current Drives:**
```go
drives, err := detector.GetRemovableDrives()
```

**DriveInfo Structure:**
```go
type DriveInfo struct {
    Name        string // Drive name/label
    Path        string // Mount point or drive letter
    Size        uint64 // Size in bytes
    SizeGB      string // Human-readable size
    IsRemovable bool   // Whether it's a removable drive
}
```

### Frontend (Svelte)

**Listening for USB Events:**
```javascript
import { EventsOn } from '../wailsjs/runtime/runtime';

EventsOn('usb-device-connected', (device) => {
    // device contains: name, path, size, sizeGB, isRemovable
});
```

## Testing

1. **Build the application:**
   ```bash
   wails build
   ```

2. **Run the application:**
   ```bash
   ./build/bin/CameraImport.exe
   ```

3. **Test USB detection:**
   - Insert a USB flash drive or SD card
   - The modal should appear within 2 seconds
   - Accept or decline the prompt

## Notes

- The detection runs automatically when the app starts
- Detection continues running in the background throughout the app's lifetime
- The polling interval (2 seconds) provides a good balance between responsiveness and system resources
- macOS implementation is event-driven (no polling) for better performance

## Future Improvements

Potential enhancements:
- Add user preference to disable auto-detection
- Remember declined devices to avoid re-prompting
- Support for auto-import when specific devices are connected
- Notification when device is disconnected during import
