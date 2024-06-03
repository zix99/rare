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

var localLinkRegexp = regexp.MustCompile(`\[[\w\s]+\]\((\w+).md\)`)

func replaceLocalLinkWithDocsCommand(match string) string {
	parts := localLinkRegexp.FindStringSubmatch(match)
	return color.Wrapf(color.Underline, "rare docs %s", parts[1])
}

var localImageRegexp = regexp.MustCompile(`!\[(.+)\]\((.+)\)`)

func replaceWithNothing(_ string) string {
	return ""
}

var regexReplacement = []regexReplacer{
	// Symbol
	{
		match:   regexp.MustCompile("`(.*?)`"),
		process: replaceWithColor(color.BrightWhite),
	},
	// Bold
	{
		match:   regexp.MustCompile(`\*\*(.*?)\*\*`),
		process: replaceWithColor(color.Bold),
	},
	// Hide images
	{
		match:   localImageRegexp,
		process: replaceWithNothing,
	},
	// Local link -> rare docs command
	{
		match:   localLinkRegexp,
		process: replaceLocalLinkWithDocsCommand,
	},
	// Raw Link
	{
		match:   regexp.MustCompile(`https?://[A-Za-z0-9.-_!#$&;=?%]+`),
		process: replaceWithColor(color.Underline),
	},
}
