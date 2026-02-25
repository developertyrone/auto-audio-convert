package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	ffmpegDirName = ".auto-audio-convert"
	ffmpegBinName = "ffmpeg"
)

// FFmpeg static build URLs (trusted sources)
var ffmpegURLs = map[string]map[string]string{
	"linux": {
		"amd64": "https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz",
		"arm64": "https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-arm64-static.tar.xz",
	},
	"darwin": {
		"amd64": "https://evermeet.cx/ffmpeg/ffmpeg-6.1.zip",
		"arm64": "https://evermeet.cx/ffmpeg/ffmpeg-6.1.zip",
	},
	"windows": {
		"amd64": "https://www.gyan.dev/ffmpeg/builds/ffmpeg-release-essentials.zip",
	},
}

// getFFmpegPath returns the path to ffmpeg (system or user-installed)
func getFFmpegPath() (string, error) {
	// Check system PATH first
	if path, err := exec.LookPath("ffmpeg"); err == nil {
		return path, nil
	}

	// Check user-installed location
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	userFFmpeg := filepath.Join(homeDir, ffmpegDirName, "bin", ffmpegBinName)
	if runtime.GOOS == "windows" {
		userFFmpeg += ".exe"
	}

	if _, err := os.Stat(userFFmpeg); err == nil {
		return userFFmpeg, nil
	}

	return "", fmt.Errorf("ffmpeg not found")
}

// promptDownloadFFmpeg asks user if they want to download ffmpeg
func promptDownloadFFmpeg() bool {
	// Check for automation flag (CI/non-interactive)
	if autoInstall := os.Getenv("AUTO_INSTALL_FFMPEG"); autoInstall == "yes" || autoInstall == "1" {
		fmt.Println("🤖 AUTO_INSTALL_FFMPEG=yes detected, downloading...")
		return true
	}

	fmt.Println("\n⚠️  FFmpeg not found!")
	fmt.Println("\nOptions:")
	fmt.Println("  1. Auto-download portable ffmpeg (~50-80MB) to ~/.auto-audio-convert/")
	fmt.Println("  2. Skip (install manually)")
	fmt.Println("\nManual installation:")
	
	switch runtime.GOOS {
	case "linux":
		fmt.Println("  sudo apt install ffmpeg       # Debian/Ubuntu")
		fmt.Println("  sudo dnf install ffmpeg       # Fedora/RHEL")
	case "darwin":
		fmt.Println("  brew install ffmpeg")
	case "windows":
		fmt.Println("  winget install ffmpeg")
		fmt.Println("  choco install ffmpeg")
	}

	fmt.Print("\nChoice [1/2]: ")
	
	var choice string
	fmt.Scanln(&choice)
	
	return choice == "1" || choice == "" || strings.ToLower(choice) == "y" || strings.ToLower(choice) == "yes"
}

// downloadFFmpeg downloads and installs portable ffmpeg
func downloadFFmpeg() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cannot get home directory: %w", err)
	}

	installDir := filepath.Join(homeDir, ffmpegDirName)
	binDir := filepath.Join(installDir, "bin")

	// Create directories
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("cannot create directory: %w", err)
	}

	// Get download URL for current platform
	osURLs, ok := ffmpegURLs[runtime.GOOS]
	if !ok {
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	downloadURL, ok := osURLs[runtime.GOARCH]
	if !ok {
		return fmt.Errorf("unsupported architecture: %s", runtime.GOARCH)
	}

	fmt.Printf("📥 Downloading ffmpeg from %s\n", getSourceName(downloadURL))
	fmt.Println("   (This may take a few minutes...)")

	// Download
	tmpFile := filepath.Join(installDir, "ffmpeg-download.tmp")
	if err := downloadFile(downloadURL, tmpFile); err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer os.Remove(tmpFile)

	fmt.Println("📦 Extracting...")

	// Extract based on file type
	if err := extractFFmpeg(tmpFile, binDir); err != nil {
		return fmt.Errorf("extraction failed: %w", err)
	}

	ffmpegPath := filepath.Join(binDir, ffmpegBinName)
	if runtime.GOOS == "windows" {
		ffmpegPath += ".exe"
	}

	// Verify installation
	if _, err := os.Stat(ffmpegPath); err != nil {
		return fmt.Errorf("ffmpeg binary not found after extraction: %w", err)
	}

	// Make executable (Unix)
	if runtime.GOOS != "windows" {
		if err := os.Chmod(ffmpegPath, 0755); err != nil {
			return fmt.Errorf("cannot make executable: %w", err)
		}
	}

	fmt.Printf("✅ FFmpeg installed to: %s\n\n", ffmpegPath)
	return nil
}

// downloadFile downloads a file from URL to destination
func downloadFile(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	// Show progress
	totalBytes := resp.ContentLength
	downloaded := int64(0)
	buffer := make([]byte, 32*1024)

	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			out.Write(buffer[:n])
			downloaded += int64(n)
			if totalBytes > 0 {
				percent := float64(downloaded) / float64(totalBytes) * 100
				fmt.Printf("\r   Progress: %.1f%% (%d/%d MB)", percent, downloaded/1024/1024, totalBytes/1024/1024)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	fmt.Println()

	return nil
}

// extractFFmpeg extracts ffmpeg binary from archive
func extractFFmpeg(archivePath, destDir string) error {
	// Detect file type by reading first few bytes
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	buf := make([]byte, 512)
	n, _ := file.Read(buf)
	file.Seek(0, 0) // Reset

	// Check magic bytes
	if n >= 4 && string(buf[:4]) == "PK\x03\x04" {
		// ZIP file
		return extractZip(archivePath, destDir)
	} else if n >= 6 && string(buf[:6]) == "\xfd7zXZ\x00" {
		// XZ compressed (tar.xz for Linux)
		return extractTarXZ(archivePath, destDir)
	} else if n >= 2 && buf[0] == 0x1f && buf[1] == 0x8b {
		// Gzip compressed (tar.gz)
		return extractTarGz(archivePath, destDir)
	}

	return fmt.Errorf("unknown archive format")
}

// extractZip extracts ffmpeg from ZIP archive
func extractZip(zipPath, destDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		baseName := filepath.Base(f.Name)
		if baseName == "ffmpeg" || baseName == "ffmpeg.exe" {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			destPath := filepath.Join(destDir, baseName)
			outFile, err := os.Create(destPath)
			if err != nil {
				return err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, rc)
			return err
		}
	}

	return fmt.Errorf("ffmpeg binary not found in archive")
}

// extractTarXZ extracts ffmpeg from tar.xz archive (Linux)
func extractTarXZ(tarPath, destDir string) error {
	// Note: tar.xz requires external xz command or a library
	// For simplicity, we'll use external tar command
	cmd := exec.Command("tar", "-xJf", tarPath, "-C", destDir, "--strip-components=1", "--wildcards", "*/ffmpeg")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, string(output))
	}
	return nil
}

// extractTarGz extracts from tar.gz (alternative for some sources)
func extractTarGz(tarPath, destDir string) error {
	file, err := os.Open(tarPath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		baseName := filepath.Base(header.Name)
		if baseName == "ffmpeg" {
			destPath := filepath.Join(destDir, baseName)
			outFile, err := os.Create(destPath)
			if err != nil {
				return err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, tr)
			return err
		}
	}

	return fmt.Errorf("ffmpeg binary not found in archive")
}

// getSourceName returns a friendly name for the download source
func getSourceName(url string) string {
	if strings.Contains(url, "johnvansickle") {
		return "johnvansickle.com (Linux static builds)"
	}
	if strings.Contains(url, "evermeet") {
		return "evermeet.cx (macOS builds)"
	}
	if strings.Contains(url, "gyan.dev") {
		return "gyan.dev (Windows builds)"
	}
	return url
}
