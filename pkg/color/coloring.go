package color

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	escapeRune     = '\x1b'
	escapeCode     = "\x1b"
	foregroundCode = "[3"
)

type ColorCode string

const (
	Reset         = escapeCode + "[0m"
	Black         = escapeCode + "[30m"
	Red           = escapeCode + "[31m"
	Green         = escapeCode + "[32m"
	Yellow        = escapeCode + "[33m"
	Blue          = escapeCode + "[34m"
	Magenta       = escapeCode + "[35m"
	Cyan          = escapeCode + "[36m"
	White         = escapeCode + "[37m"
	BrightBlack   = escapeCode + "[30;1m"
	BrightRed     = escapeCode + "[31;1m"
	BrightGreen   = escapeCode + "[32;1m"
	BrightYellow  = escapeCode + "[33;1m"
	BrightBlue    = escapeCode + "[34;1m"
	BrightMagenta = escapeCode + "[35;1m"
	BrightCyan    = escapeCode + "[36;1m"
	BrightWhite   = escapeCode + "[37;1m"

	Bold      = escapeCode + "[1m"
	Underline = escapeCode + "[4m"
)

var colorMap = map[string]ColorCode{
	"black":   Black,
	"red":     Red,
	"green":   Green,
	"yellow":  Yellow,
	"blue":    Blue,
	"magenta": Magenta,
	"cyan":    Cyan,
	"white":   White,
}

// Enabled controls whether or not coloring is applied
var Enabled = true

var GroupColors = [...]ColorCode{Red, Green, Yellow, Blue, Magenta, Cyan, BrightRed, BrightGreen, BrightYellow, BrightBlue, BrightMagenta, BrightCyan}

func init() {
	if fi, err := os.Stdout.Stat(); err == nil {
		if (fi.Mode() & os.ModeCharDevice) == 0 {
			Enabled = false
		}
	}
}

func Write(w io.StringWriter, color ColorCode, f func(w io.StringWriter)) {
	if !Enabled {
		f(w)
		return
	}

	w.WriteString(string(color))
	f(w)
	w.WriteString(string(Reset))
}

// Wrap surroungs a string with a color (if enabled)
func Wrap(color ColorCode, s string) string {
	if !Enabled {
		return s
	}

	var sb strings.Builder
	sb.Grow(len(s) + 8)
	sb.WriteString(string(color))
	sb.WriteString(s)
	if len(s) < len(Reset) || s[len(s)-len(Reset):] != string(Reset) {
		sb.WriteString(string(Reset))
	}
	return sb.String()
}

func Wrapf(color ColorCode, s string, args ...interface{}) string {
	return Wrap(color, fmt.Sprintf(s, args...))
}

func Wrapi(color ColorCode, s interface{}) string {
	return Wrap(color, fmt.Sprintf("%v", s))
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
		if start >= 0 && end >= 0 && end > start && start >= lastIndex {
			color := GroupColors[(i/2)%len(GroupColors)]

			sb.WriteString(s[lastIndex:start])
			sb.WriteString(string(color))
			sb.WriteString(s[start:end])
			sb.WriteString(string(Reset))

			lastIndex = end
		}
	}

	if lastIndex < len(s) {
		sb.WriteString(s[lastIndex:])
	}

	return sb.String()
}

func LookupColorByName(s string) (ColorCode, bool) {
	if c, ok := colorMap[strings.ToLower(s)]; ok {
		return c, true
	}
	return BrightRed, false
}

// StrLen ignoring any color codes. If color disabled, returns len(s)
func StrLen(s string) (ret int) {
	if !Enabled {
		return len(s)
	}

	inCode := false
	for _, r := range s {
		if r == escapeRune {
			inCode = true
		} else if inCode && r == 'm' {
			inCode = false
		} else if !inCode {
			ret++
		}
	}
	return
}
