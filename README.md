# Auto Audio Convert

A lightweight, cross-platform batch audio converter that automatically scans directories and converts audio files using FFmpeg.

## Features

- 🚀 **Efficient**: Worker pool with configurable concurrency
- 🔄 **Recursive**: Scans subdirectories automatically
- ⏭️ **Smart Skip**: Avoids re-converting existing files
- 💾 **Resource-Limited**: 512MB memory cap, CPU throttling
- 🌍 **Cross-Platform**: Single binary for Windows, macOS, Linux
- 📦 **Zero Dependencies**: Just Go stdlib + FFmpeg

## Prerequisites

**FFmpeg** is required for audio conversion.

### Auto-Install (Recommended)
The first time you run `auto-audio-convert`, if FFmpeg is not found, you'll be prompted to auto-download a portable version:

```bash
./auto-audio-convert --from=flac --to=mp3

⚠️  FFmpeg not found!

Options:
  1. Auto-download portable ffmpeg (~50-80MB) to ~/.auto-audio-convert/
  2. Skip (install manually)
```

Choose option **1** for automatic installation (downloads to your home directory, no sudo required).

**For CI/automation:** Set `AUTO_INSTALL_FFMPEG=yes` to skip the prompt.

### Manual Install
- Ubuntu/Debian: `sudo apt install ffmpeg`
- macOS: `brew install ffmpeg`
- Windows: `winget install ffmpeg` or download from [ffmpeg.org](https://ffmpeg.org/download.html)

## Installation

### Download Pre-built Binary (Coming Soon)
```bash
# Linux/macOS
curl -L https://github.com/developertyrone/auto-audio-convert/releases/latest/download/auto-audio-convert-$(uname -s)-$(uname -m) -o auto-audio-convert
chmod +x auto-audio-convert
```

### Build from Source
```bash
git clone https://github.com/developertyrone/auto-audio-convert.git
cd auto-audio-convert
go build -o auto-audio-convert
```

## Usage

### Basic Conversion
```bash
# Convert all FLAC files to MP3 in current directory
./auto-audio-convert --from=flac --to=mp3

# Convert WAV to OGG in specific directory
./auto-audio-convert --source=/path/to/music --from=wav --to=ogg

# Use 8 parallel workers
./auto-audio-convert --from=flac --to=mp3 --workers=8
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--source` | Source directory to scan | `.` (current dir) |
| `--from` | Source file extension | **(required)** |
| `--to` | Target file extension | **(required)** |
| `--workers` | Number of parallel workers | CPU cores / 2 |
| `--version` | Show version | - |

## Resource Limits

- **Memory**: 512MB soft limit (Go runtime)
- **CPU**: Default workers = `NumCPU/2`
- **FFmpeg**: 2 threads per conversion

Adjust `--workers` based on your system:
- **Laptop/Desktop**: 2-4 workers
- **Server**: 8-16 workers

## Build for Multiple Platforms

```bash
# Linux (64-bit)
GOOS=linux GOARCH=amd64 go build -o auto-audio-convert-linux-amd64

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o auto-audio-convert-darwin-arm64

# Windows (64-bit)
GOOS=windows GOARCH=amd64 go build -o auto-audio-convert-windows-amd64.exe
```

## Example Output

```
🎵 Auto Audio Converter v1.0.0
📁 Source: /home/user/music
🔄 Converting: .flac → .mp3
⚙️  Workers: 4 (CPU limit: 4 cores)
💾 Memory limit: 512MB

📂 Found 12 .flac file(s)

🔄 [Worker 0] Converting: song1.flac
✅ [Worker 0] Done: song1.flac
⏭️  [Worker 1] Skipped: song2.flac (target exists: song2.mp3)
🔄 [Worker 2] Converting: song3.flac
✅ [Worker 2] Done: song3.flac

📊 Summary:
   ✅ Converted: 10
   ⏭️  Skipped: 2
   ❌ Failed: 0
```

## License

MIT

## Contributing

PRs welcome! Please ensure:
1. Code passes `go fmt` and `go vet`
2. Add tests for new features
3. Update README for new flags/features
