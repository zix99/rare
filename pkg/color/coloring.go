package color

import (
	"strings"
)

const (
	escapeCode     = "\x1b"
	foregroundCode = "[3"
)

type ColorCode string

const (
	Reset   ColorCode = escapeCode + "[0m"
	Red               = escapeCode + "[31m"
	Green             = escapeCode + "[32m"
	Yellow            = escapeCode + "[33m"
	Blue              = escapeCode + "[34m"
	Magenta           = escapeCode + "[35m"
	Cyan              = escapeCode + "[36m"
)

// Enabled controls whether or not coloring is applied
var Enabled = true

var groupColors = [...]ColorCode{Red, Green, Yellow, Blue, Magenta, Cyan}

// Wrap surroungs a string with a color (if enabled)
func Wrap(color ColorCode, s string) string {
	if !Enabled {
		return s
	}

	var sb strings.Builder
	sb.Grow(len(s) + 8)
	sb.WriteString(string(color))
	sb.WriteString(s)
	sb.WriteString(string(Reset))
	return sb.String()
}

// WrapIndices color-codes by group pairs (regex-style)
//  [aStart, aEnd, bStart, bEnd...]
func WrapIndices(s string, groups []int) string {
	if !Enabled {
		return s
	}
	if len(groups) == 0 || len(groups)%2 != 0 {
		return s
	}

	var sb strings.Builder
	lastIndex := 0

	for i := 0; i < len(groups); i += 2 {
		start := groups[i]
		end := groups[i+1]
		color := groupColors[(i/2)%len(groupColors)]

		sb.WriteString(s[lastIndex:start])
		sb.WriteString(string(color))
		sb.WriteString(s[start:end])
		sb.WriteString(string(Reset))

		lastIndex = end
	}

	if lastIndex < len(s) {
		sb.WriteString(s[lastIndex:])
	}

	return sb.String()
}
