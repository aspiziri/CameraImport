# Troubleshooting Guide

## Binding Generation Errors

### Error: "cannot embed irregular files"
This happens when the frontend/dist directory doesn't exist.

**Fix:**
```bash
mkdir -p frontend/dist
echo "placeholder" > frontend/dist/.gitkeep
```

### Error: "cannot generate bindings for unexported type"
All types in method signatures must start with a capital letter (exported).

**Fix:** All our types are already exported, but if you see this error, check that custom types used in method parameters or return values start with capital letters.

### Error: "wails.json: no such file or directory"
Wails can't find the config file.

**Fix:** Make sure you're running `wails dev` from the project root directory.

### Error: "frontend:install script failed"
The npm install failed.

**Fix:**
```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
cd ..
```

### Error: "go: cannot find module"
Go dependencies not installed.

**Fix:**
```bash
go mod download
go mod tidy
```

## Testing Compilation Without Wails

### Test Go Backend Only
```bash
# In PowerShell or CMD
cd internal/importer
go build
cd ../..
```

### Test Frontend Only
```bash
cd frontend
npm run build
cd ..
```

## Common Issues

### 1. FFmpeg/FFprobe not found
The video metadata extraction requires ffprobe.

**Fix:**
- Install ffmpeg from https://ffmpeg.org/download.html
- Add to PATH
- Verify: `ffprobe -version`

### 2. MinGW not found (Windows)
Wails requires GCC on Windows.

**Fix:**
- Install MSYS2 from https://www.msys2.org/
- Run: `pacman -S mingw-w64-x86_64-gcc`
- Add `C:\msys64\mingw64\bin` to PATH

### 3. WebView2 not found (Windows)
Wails needs WebView2 runtime.

**Fix:**
- Usually pre-installed on Windows 11
- Download from: https://developer.microsoft.com/en-us/microsoft-edge/webview2/

## Manual Binding Generation

If automatic binding generation fails, you can try generating them manually:

```bash
wails generate module
```

## Simplified Test

To test if the basic structure works, you can temporarily simplify the ImportService:

1. Comment out complex methods in `internal/importer/importer.go`
2. Keep only simple methods like:
   ```go
   func (s *ImportService) Test() string {
       return "Hello from Go!"
   }
   ```
3. Run `wails dev` to verify bindings work
4. Gradually uncomment methods to find the problematic one

## Getting More Info

Run wails with verbose output:
```bash
wails dev -v
```

Or build with debug info:
```bash
wails build -debug
```

## Still Having Issues?

Please provide:
1. Full error message from terminal
2. Output of `wails doctor`
3. Your OS and versions:
   - `go version`
   - `node --version`
   - `wails version`
