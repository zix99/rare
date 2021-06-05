package patterns

import (
	"embed"
	"path/filepath"
	"rare/pkg/logger"
)

//go:embed data/*
var patternFiles embed.FS

const patternPath = "data"

var stdLib *PatternSet

func Stdlib() *PatternSet {
	if stdLib == nil {
		stdLib = NewPatternSet()
		entries, _ := patternFiles.ReadDir(patternPath)
		logger.Print(entries)
		for _, entry := range entries {
			fullPath := filepath.Join(patternPath, entry.Name())
			f, _ := patternFiles.Open(fullPath)
			defer f.Close()

			stdLib.LoadPatternFile(f)
		}
	}

	return stdLib
}

// Lookup is ease-of-access to stdlib().lookup()
func Lookup(p string) (string, bool) {
	return Stdlib().Lookup(p)
}
