package main

import (
	"io/fs"
	"path/filepath"
	"strings"
)

// scanDirectory walks the directory tree and returns all files with the target extension
func scanDirectory(root string, ext string) ([]string, error) {
	var files []string
	ext = strings.ToLower(strings.TrimPrefix(ext, "."))

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		fileExt := strings.ToLower(strings.TrimPrefix(filepath.Ext(path), "."))
		if fileExt == ext {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}
