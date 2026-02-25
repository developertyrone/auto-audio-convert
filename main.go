package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/debug"
)

const (
	defaultMemLimitMB = 512
	version           = "1.0.0" // Can be overridden by -ldflags at build time
)

func main() {
	// Set memory limit (Go 1.19+)
	debug.SetMemoryLimit(defaultMemLimitMB * 1024 * 1024)

	// CLI flags
	sourceDir := flag.String("source", ".", "Source directory to scan")
	sourceExt := flag.String("from", "", "Source file extension (e.g., wav, flac)")
	targetExt := flag.String("to", "", "Target file extension (e.g., mp3, ogg)")
	workers := flag.Int("workers", runtime.NumCPU()/2, "Number of parallel workers")
	showVersion := flag.Bool("version", false, "Show version")

	flag.Parse()

	if *showVersion {
		fmt.Printf("auto-audio-convert v%s\n", version)
		os.Exit(0)
	}

	// Validate inputs
	if *sourceExt == "" || *targetExt == "" {
		log.Fatal("Error: Both --from and --to extensions are required\n\nUsage:\n  auto-audio-convert --source=/path/to/audio --from=flac --to=mp3 [--workers=4]\n")
	}

	if *workers < 1 {
		*workers = 1
	}

	fmt.Printf("🎵 Auto Audio Converter v%s\n", version)
	fmt.Printf("📁 Source: %s\n", *sourceDir)
	fmt.Printf("🔄 Converting: .%s → .%s\n", *sourceExt, *targetExt)
	fmt.Printf("⚙️  Workers: %d (CPU limit: %d cores)\n", *workers, runtime.NumCPU()/2)
	fmt.Printf("💾 Memory limit: %dMB\n\n", defaultMemLimitMB)

	// Check ffmpeg availability
	ffmpegPath, err := getFFmpegPath()
	if err != nil {
		// Prompt to download
		if promptDownloadFFmpeg() {
			if err := downloadFFmpeg(); err != nil {
				log.Fatalf("❌ Failed to download ffmpeg: %v\n", err)
			}
			// Re-check after download
			ffmpegPath, err = getFFmpegPath()
			if err != nil {
				log.Fatal("❌ FFmpeg installation failed. Please install manually.")
			}
		} else {
			log.Fatal("❌ FFmpeg required. Please install it and try again.")
		}
	}

	fmt.Printf("🎬 Using ffmpeg: %s\n", ffmpegPath)

	// Scan and convert
	files, err := scanDirectory(*sourceDir, *sourceExt)
	if err != nil {
		log.Fatalf("❌ Scan error: %v", err)
	}

	if len(files) == 0 {
		fmt.Printf("✅ No .%s files found in %s\n", *sourceExt, *sourceDir)
		return
	}

	fmt.Printf("📂 Found %d .%s file(s)\n\n", len(files), *sourceExt)

	// Process with worker pool
	stats := convertFiles(files, *targetExt, *workers, ffmpegPath)

	// Summary
	fmt.Printf("\n📊 Summary:\n")
	fmt.Printf("   ✅ Converted: %d\n", stats.Converted)
	fmt.Printf("   ⏭️  Skipped: %d\n", stats.Skipped)
	fmt.Printf("   ❌ Failed: %d\n", stats.Failed)
}
