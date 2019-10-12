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

var Enabled = true

var groupColors = [...]ColorCode{Red, Green, Yellow, Blue, Magenta, Cyan}

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

func ColorCodeGroups(s string, groups []string) string {
	if !Enabled {
		return s
	}
	return s
}
