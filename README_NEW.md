# CameraImport v4.0 - Desktop Application

A modern desktop application for importing and organizing photos and videos from cameras and SD cards.

![CameraImport](https://img.shields.io/badge/version-4.0.0-blue)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![Wails](https://img.shields.io/badge/Wails-v2-ff69b4)
![Svelte](https://img.shields.io/badge/Svelte-4-FF3E00?logo=svelte)

## Overview

CameraImport v4.0 is a complete rewrite of the Python script as a modern desktop application using:
- **Backend**: Go (fast, efficient, single binary)
- **Frontend**: Svelte (reactive, lightweight UI)
- **Framework**: Wails v2 (Electron alternative, smaller footprint)

## Key Features

### 🚀 Core Functionality
- **Full Import**: One-click import of all media files
- **Selective Import**: Choose specific dates and custom folder names
- **Smart Organization**: Automatic date-based folder structure
- **Duplicate Detection**: Skip files that already exist

### 📸 Advanced Features
- **EXIF Reading**: Uses actual photo capture date from EXIF metadata
- **Live Preview**: Thumbnail previews of images before import
- **Video Metadata**: Extracts resolution, FPS, and duration for videos
- **Progress Tracking**: Real-time import progress with file counts
- **Configurable Formats**: Support for any image, video, or RAW format

### 🎨 Modern UI
- Dark theme with smooth animations
- Responsive layout
- Drag-and-drop folder selection (planned)
- File grid and date list views
- Import profiles/presets (planned)

## Quick Start

### For Users (Running the app)

1. Download the latest release for your platform
2. Run the executable
3. Configure source (SD card) and destination folders
4. Click "Scan Files" to preview
5. Choose Full or Selective import
6. Click "Start Import"

### For Developers (Building from source)

See [BUILD.md](BUILD.md) for detailed build instructions.

**Quick build:**
```bash
# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Install dependencies
go mod download
cd frontend && npm install && cd ..

# Run in dev mode
wails dev

# Or build for production
wails build
```

## File Organization

### Default Structure
```
Destination/
├── 2024-01-15/
│   ├── Capture/
│   │   ├── 2024-01-15 14.23.45.ARW
│   │   └── 2024-01-15 14.24.12.ARW
│   ├── Videos/
│   │   └── 2024-01-15 15.30.22 - 3840w - 29.970 - 45sec.MP4
│   ├── 2024-01-15 14.23.45.JPG
│   └── 2024-01-15 14.24.12.JPG
└── 2024-01-16/
    └── ...
```

### Filename Formats
- **Images**: `YYYY-MM-DD HH.MM.SS.jpg`
- **RAW**: `YYYY-MM-DD HH.MM.SS.arw`
- **Videos**: `YYYY-MM-DD HH.MM.SS - RESw - FPS - DURsec.mp4`

## Configuration

### Supported Formats
Configure which file formats to import:

**Default supported formats:**
- **Images**: JPG, JPEG, PNG, GIF, TIF
- **Videos**: MP4, AVI, MOV
- **RAW**: ARW, CR2, NEF, DNG

### Folder Structure
Customize where each file type goes:
- `imgRelativePath`: Where images are stored (default: `/`)
- `videoRelativePath`: Where videos are stored (default: `/Videos/`)
- `rawRelativePath`: Where RAW files are stored (default: `/Capture/`)

## Comparison: v3 (Python) vs v4 (Wails)

| Feature | v3.0 (Python) | v4.0 (Wails) |
|---------|---------------|--------------|
| UI | Terminal CLI | Modern Desktop App |
| File Preview | No | Yes (thumbnails) |
| EXIF Reading | No | Yes |
| Video Metadata | Yes (ffprobe) | Yes (ffprobe) |
| Platform | Python required | Standalone binary |
| Size | ~5KB (script) | ~10-15MB (app) |
| Performance | Good | Excellent |
| Dev Stack | Python | Go + Svelte |

## System Requirements

### Runtime (Users)
- **Windows**: Windows 10/11 (WebView2 included)
- **macOS**: macOS 10.15 or later
- **Linux**: GTK3 and WebKit2GTK
- **FFmpeg**: Required for video metadata

### Build (Developers)
- Go 1.21+
- Node.js 18+
- Wails CLI
- Platform-specific tools (see BUILD.md)

## Roadmap

- [x] Full import mode
- [x] Selective import mode
- [x] EXIF data reading
- [x] Live file preview
- [x] Video metadata extraction
- [ ] Drag-and-drop folder selection
- [ ] Import profiles/presets
- [ ] Batch image editing (rotate, resize)
- [ ] Watermark support
- [ ] Cloud storage integration
- [ ] Import scheduling

## Contributing

Contributions are welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## Migration from v3 (Python)

Your existing `config.ini` is not directly compatible. The new app uses JSON configuration, but you can:
1. Run the new app
2. Configure settings in the UI
3. Settings are saved automatically

Your folder structure and organization preferences remain the same.

## License

MIT License - See LICENSE file for details.

## Credits

Built with:
- [Wails](https://wails.io/) - Desktop app framework
- [Svelte](https://svelte.dev/) - Frontend framework
- [Go](https://golang.org/) - Backend language
- [goexif](https://github.com/rwcarlsen/goexif) - EXIF reading
- [imaging](https://github.com/disintegration/imaging) - Image processing

---

**Previous version**: The Python v3.0 script is still available in `CameraImport.py` for legacy use.
