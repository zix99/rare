package stdpat

import (
	"embed"
	"rare/pkg/grok/patterns"
	"rare/pkg/logger"
)

//go:embed *.grok
var patternFiles embed.FS

var stdLib *patterns.PatternSet

func Stdlib() *patterns.PatternSet {
	if stdLib == nil {
		stdLib = patterns.NewPatternSet()
		entries, _ := patternFiles.ReadDir(".")
		logger.Print(entries)
		for _, entry := range entries {
			f, _ := patternFiles.Open(entry.Name())
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
