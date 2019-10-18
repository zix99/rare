package markdowncli

import (
	"rare/pkg/color"
	"regexp"
)

type regexReplacerFunc func(string) string

type regexReplacer struct {
	match   *regexp.Regexp
	process regexReplacerFunc
}

func replaceWithColor(clr color.ColorCode) regexReplacerFunc {
	return regexReplacerFunc(func(match string) string {
		return color.Wrap(clr, match)
	})
}

var regexReplacement = []regexReplacer{
	// Symbol
	regexReplacer{
		match:   regexp.MustCompile("`(.*?)`"),
		process: replaceWithColor(color.BrightWhite),
	},
	// Bold
	regexReplacer{
		match:   regexp.MustCompile(`\*\*(.*?)\*\*`),
		process: replaceWithColor(color.Bold),
	},
	// Raw Link
	regexReplacer{
		match:   regexp.MustCompile(`https?://[A-Za-z0-9.-_!#$&;=?%]+`),
		process: replaceWithColor(color.Underline),
	},
}
