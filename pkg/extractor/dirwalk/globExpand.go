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
