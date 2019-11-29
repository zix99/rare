package helpers

import (
	"os"
	"path/filepath"
)

func isDir(path string) bool {
	if fi, err := os.Stat(path); err == nil && fi.IsDir() {
		return true
	}
	return false
}

// globExpand expands a directory-glob, and optionally recurses on it
func globExpand(paths []string, recursive bool) []string {
	out := make([]string, 0)
	for _, p := range paths {
		if recursive && isDir(p) {
			filepath.Walk(p, func(walkPath string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					out = append(out, walkPath)
				}
				return nil
			})
		} else {
			expanded, err := filepath.Glob(p)
			if err != nil {
				ErrLog.Printf("Path error: %v\n", err)
			} else {
				out = append(out, expanded...)
			}
		}
	}
	return out
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
