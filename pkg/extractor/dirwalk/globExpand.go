package dirwalk

import (
	"os"
	"path/filepath"
	"rare/pkg/logger"
)

// globExpand expands a directory-glob, and optionally recurses on it asynchronously
func GlobExpand(paths []string, recursive bool) <-chan string {
	c := make(chan string, 10)
	go func() {
		for _, p := range paths {
			if recursive && isDir(p) {
				filepath.WalkDir(p, func(walkPath string, info os.DirEntry, err error) error {
					if err != nil {
						logger.Printf("Path error: %v", err)
					} else if info.Type().IsRegular() {
						c <- walkPath
					}
					return nil
				})
			} else {
				expanded, err := filepath.Glob(p)
				if err != nil {
					logger.Printf("Path error: %v", err)
				} else if len(expanded) > 0 {
					for _, item := range expanded {
						c <- item
					}
				} else {
					c <- p
				}
			}
		}
		close(c)
	}()
	return c
}

func isDir(path string) bool {
	if fi, err := os.Stat(path); err == nil && fi.IsDir() {
		return true
	}
	return false
}
