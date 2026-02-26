package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
)

// ConversionStats tracks conversion results
type ConversionStats struct {
	Converted int32
	Skipped   int32
	Failed    int32
}

// targetFileExists checks if the target file already exists
func targetFileExists(sourcePath, targetExt string) (bool, string) {
	ext := filepath.Ext(sourcePath)
	targetPath := strings.TrimSuffix(sourcePath, ext) + "." + strings.TrimPrefix(targetExt, ".")

	if _, err := os.Stat(targetPath); err == nil {
		return true, targetPath
	}

	return false, targetPath
}

// convertFile converts a single audio file using ffmpeg
func convertFile(sourcePath, targetExt, ffmpegPath, bitrate string, overwrite bool) error {
	exists, targetPath := targetFileExists(sourcePath, targetExt)
	if exists && !overwrite {
		return fmt.Errorf("target already exists")
	}

	// Build FFmpeg command with resource limits
	args := []string{
		"-i", sourcePath,
		"-threads", "2",
		"-loglevel", "error",
	}

	// Add bitrate if specified
	if bitrate != "" {
		args = append(args, "-b:a", bitrate)
	}

	// Add output file
	args = append(args, "-y", targetPath)

	cmd := exec.Command(ffmpegPath, args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, string(output))
	}

	return nil
}

// convertFiles processes files using a worker pool
func convertFiles(files []string, targetExt string, workers int, ffmpegPath, bitrate string, overwrite bool) ConversionStats {
	var stats ConversionStats
	var wg sync.WaitGroup

	// Create job channel
	jobs := make(chan string, len(files))

	// Start workers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for sourcePath := range jobs {
				relPath := filepath.Base(sourcePath)

				// Check if already converted (unless overwrite is enabled)
				exists, targetPath := targetFileExists(sourcePath, targetExt)
				if exists && !overwrite {
					fmt.Printf("⏭️  [Worker %d] Skipped: %s (target exists: %s)\n",
						workerID, relPath, filepath.Base(targetPath))
					atomic.AddInt32(&stats.Skipped, 1)
					continue
				}

				// Convert
				if exists && overwrite {
					fmt.Printf("🔄 [Worker %d] Overwriting: %s\n", workerID, relPath)
				} else {
					fmt.Printf("🔄 [Worker %d] Converting: %s\n", workerID, relPath)
				}
				err := convertFile(sourcePath, targetExt, ffmpegPath, bitrate, overwrite)

				if err != nil {
					if strings.Contains(err.Error(), "target already exists") {
						fmt.Printf("⏭️  [Worker %d] Skipped: %s (target exists)\n", workerID, relPath)
						atomic.AddInt32(&stats.Skipped, 1)
					} else {
						fmt.Printf("❌ [Worker %d] Failed: %s - %v\n", workerID, relPath, err)
						atomic.AddInt32(&stats.Failed, 1)
					}
				} else {
					fmt.Printf("✅ [Worker %d] Done: %s\n", workerID, relPath)
					atomic.AddInt32(&stats.Converted, 1)
				}
			}
		}(i)
	}

	// Send jobs
	for _, file := range files {
		jobs <- file
	}
	close(jobs)

	// Wait for completion
	wg.Wait()

	return stats
}
