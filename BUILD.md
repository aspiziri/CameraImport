# Building CameraImport

This guide will help you build and run the CameraImport desktop application.

## Prerequisites

### Required Software

1. **Go** (version 1.21 or higher)
   - Download: https://golang.org/dl/
   - Verify: `go version`

2. **Node.js** (version 18 or higher)
   - Download: https://nodejs.org/
   - Verify: `node --version` and `npm --version`

3. **Wails CLI**
   ```bash
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   ```
   - Verify: `wails version`

4. **FFmpeg** (for video metadata extraction)
   - Windows: Download from https://ffmpeg.org/download.html
   - Add to PATH
   - Verify: `ffprobe -version`

5. **GCC/MinGW** (Windows only)
   - Download: https://www.msys2.org/
   - After installing MSYS2, run:
     ```bash
     pacman -S mingw-w64-x86_64-gcc
     ```
   - Add `C:\msys64\mingw64\bin` to PATH

### Platform-Specific Requirements

#### Windows
- Visual Studio Build Tools or Visual Studio Community
- WebView2 Runtime (usually pre-installed on Windows 11)

#### macOS
- Xcode Command Line Tools: `xcode-select --install`

#### Linux
- GTK3 and WebKit2GTK
  ```bash
  # Ubuntu/Debian
  sudo apt install libgtk-3-dev libwebkit2gtk-4.0-dev

  # Fedora
  sudo dnf install gtk3-devel webkit2gtk3-devel
  ```

## Installation

1. **Clone or navigate to the project directory**
   ```bash
   cd CameraImport
   ```

2. **Install Go dependencies**
   ```bash
   go mod download
   ```

3. **Install frontend dependencies**
   ```bash
   cd frontend
   npm install
   cd ..
   ```

## Development

### Run in development mode
```bash
wails dev
```

This will:
- Start the Go backend
- Start the Vite dev server for hot-reloading
- Open the application window
- Enable browser devtools

## Building

### Build for your current platform
```bash
wails build
```

The built application will be in the `build/bin/` directory.

### Build with optimizations (smaller binary)
```bash
wails build -clean -ldflags "-w -s"
```

### Cross-platform builds

#### Build for Windows (from any platform)
```bash
wails build -platform windows/amd64
```

#### Build for macOS (requires macOS)
```bash
wails build -platform darwin/universal
```

#### Build for Linux
```bash
wails build -platform linux/amd64
```

## Project Structure

```
CameraImport/
├── main.go                    # Application entry point
├── app.go                     # App struct and methods
├── internal/                  # Internal Go packages
│   ├── config/               # Configuration management
│   ├── exif/                 # EXIF data reading
│   ├── importer/             # Core import logic
│   └── preview/              # Thumbnail generation
├── frontend/                  # Frontend Svelte application
│   ├── src/
│   │   ├── App.svelte       # Main app component
│   │   ├── components/      # UI components
│   │   └── main.js          # Frontend entry point
│   ├── package.json
│   └── vite.config.js
└── build/                     # Build output (generated)
```

## Configuration

The application uses JSON configuration files. You can create a `config.json` in the same directory as the executable:

```json
{
  "imgFormats": ["jpg", "JPG", "png", "PNG"],
  "videoFormats": ["mp4", "MP4", "mov", "MOV"],
  "rawFormats": ["arw", "ARW", "cr2", "CR2"],
  "source": "",
  "destination": "",
  "imgRelativePath": "/",
  "videoRelativePath": "/Videos/",
  "rawRelativePath": "/Capture/"
}
```

## Features

### Core Features
- ✅ Full and selective import modes
- ✅ EXIF data reading for accurate capture dates
- ✅ Video metadata extraction (resolution, FPS, duration)
- ✅ Live file preview with thumbnails
- ✅ Date-based organization
- ✅ Duplicate detection
- ✅ Custom folder naming
- ✅ Progress tracking
- ✅ Optional file deletion after import

### File Organization
Files are organized into date-based folders with configurable subfolders:
- Images: `/YYYY-MM-DD/` (configurable)
- Videos: `/YYYY-MM-DD/Videos/` (configurable)
- RAW: `/YYYY-MM-DD/Capture/` (configurable)

### Filename Format
- **Images & RAW**: `YYYY-MM-DD HH.MM.SS.ext`
- **Videos**: `YYYY-MM-DD HH.MM.SS - WIDTHw - FPS - DURATIONsec.ext`

## Troubleshooting

### "wails: command not found"
Make sure `$GOPATH/bin` is in your PATH:
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### "gcc: command not found" (Windows)
Install MinGW-w64 via MSYS2 and add to PATH.

### Frontend build errors
```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
cd ..
wails build
```

### FFprobe not working
Ensure ffmpeg is installed and `ffprobe` is in your PATH:
```bash
ffprobe -version
```

## License

MIT License - See LICENSE file for details.

## Support

For issues and questions, please open an issue on the project repository.
