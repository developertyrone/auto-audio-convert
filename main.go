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
	version           = "1.2.0" // Can be overridden by -ldflags at build time
)

// resolveBitrate converts quality preset or validates custom bitrate
func resolveBitrate(quality, bitrate string) string {
	// Custom bitrate takes priority
	if bitrate != "" {
		return bitrate
	}

	// Quality presets
	switch quality {
	case "low":
		return "128k"
	case "medium":
		return "192k"
	case "high":
		return "320k"
	case "":
		return "" // No quality specified - use ffmpeg defaults
	default:
		log.Printf("⚠️  Warning: Unknown quality preset '%s', using ffmpeg defaults\n", quality)
		return ""
	}
}

func main() {
	// Set memory limit (Go 1.19+)
	debug.SetMemoryLimit(defaultMemLimitMB * 1024 * 1024)

	// CLI flags
	sourceDir := flag.String("source", ".", "Source directory to scan")
	sourceExt := flag.String("from", "", "Source file extension (e.g., wav, flac)")
	targetExt := flag.String("to", "", "Target file extension (e.g., mp3, ogg)")
	workers := flag.Int("workers", runtime.NumCPU()/2, "Number of parallel workers")
	quality := flag.String("quality", "", "Quality preset: low (128k), medium (192k), high (320k)")
	bitrate := flag.String("bitrate", "", "Custom bitrate (e.g., 256k, 320k) - overrides quality preset")
	overwrite := flag.Bool("overwrite", false, "Overwrite existing target files (skip existence check)")
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

	// Resolve quality/bitrate
	audioBitrate := resolveBitrate(*quality, *bitrate)

	fmt.Printf("🎵 Auto Audio Converter v%s\n", version)
	fmt.Printf("📁 Source: %s\n", *sourceDir)
	fmt.Printf("🔄 Converting: .%s → .%s\n", *sourceExt, *targetExt)
	if audioBitrate != "" {
		fmt.Printf("🎚️  Quality: %s\n", audioBitrate)
	}
	if *overwrite {
		fmt.Printf("⚠️  Mode: OVERWRITE (existing files will be replaced)\n")
	}
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
	stats := convertFiles(files, *targetExt, *workers, ffmpegPath, audioBitrate, *overwrite)

	// Summary
	fmt.Printf("\n📊 Summary:\n")
	fmt.Printf("   ✅ Converted: %d\n", stats.Converted)
	fmt.Printf("   ⏭️  Skipped: %d\n", stats.Skipped)
	fmt.Printf("   ❌ Failed: %d\n", stats.Failed)
}
