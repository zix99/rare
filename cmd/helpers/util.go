package helpers

import (
	"os"
	"path/filepath"
	"rare/pkg/logger"
)

// Aggregate one channel into another, with a buffer
func bufferChan(in <-chan string, size int) <-chan string {
	out := make(chan string, size)
	go func() {
		for item := range in {
			out <- item
		}
		close(out)
	}()
	return out
}

func isDir(path string) bool {
	if fi, err := os.Stat(path); err == nil && fi.IsDir() {
		return true
	}
	return false
}

// globExpand expands a directory-glob, and optionally recurses on it
func globExpand(paths []string, recursive bool) <-chan string {
	c := make(chan string, 10)
	go func() {
		for _, p := range paths {
			if recursive && isDir(p) {
				filepath.Walk(p, func(walkPath string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if !info.IsDir() {
						c <- walkPath
					}
					return nil
				})
			} else {
				expanded, err := filepath.Glob(p)
				if err != nil {
					logger.Printf("Path error: %v", err)
				} else {
					for _, item := range expanded {
						c <- item
					}
				}
			}
		}
		close(c)
	}()
	return c
}
