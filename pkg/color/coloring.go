package color

import (
	"fmt"
	"io"
	"rare/pkg/multiterm/termstate"
	"strconv"
	"strings"
	"unicode/utf8"
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
	if termstate.IsPipedOutput() {
		Enabled = false
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

func WriteString(w io.StringWriter, color ColorCode, s string) {
	if !Enabled {
		w.WriteString(s)
		return
	}

	w.WriteString(string(color))
	w.WriteString(s)
	w.WriteString(Reset)
}

func WriteUint64(w io.StringWriter, color ColorCode, v uint64) {
	sv := strconv.FormatUint(v, 10)

	if !Enabled {
		w.WriteString(sv)
		return
	}

	w.WriteString(string(color))
	w.WriteString(sv)
	w.WriteString(Reset)
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

func Wrapi(color ColorCode, v int) string {
	return Wrap(color, strconv.Itoa(v))
}

// WrapIndices color-codes by group pairs (regex-style)
//
//	[aStart, aEnd, bStart, bEnd...]
func WrapIndices(sw io.StringWriter, s string, groups []int) {
	if !Enabled {
		sw.WriteString(s)
		return
	}
	if len(groups) == 0 || len(groups)%2 != 0 {
		sw.WriteString(s)
		return
	}
	lastIndex := 0

	for i := 0; i < len(groups); i += 2 {
		start := groups[i]
		end := groups[i+1]
		if start >= 0 && end >= 0 && end > start && start >= lastIndex {
			color := GroupColors[(i/2)%len(GroupColors)]

			sw.WriteString(s[lastIndex:start])
			sw.WriteString(string(color))
			sw.WriteString(s[start:end])
			sw.WriteString(string(Reset))

			lastIndex = end
		}
	}

	if lastIndex < len(s) {
		sw.WriteString(s[lastIndex:])
	}
}

func LookupColorByName(s string) (ColorCode, bool) {
	if c, ok := colorMap[strings.ToLower(s)]; ok {
		return c, true
	}
	return BrightRed, false
}

// UnderlineSingleRune is a special-use-case colorer for headers with a single character called out
func WriteHighlightSingleRune(w io.StringWriter, word string, runeIndex int, base, highlight ColorCode) {
	if !Enabled {
		w.WriteString(word)
		return
	}

	w.WriteString(string(base))

	offset, width := byteIndexOfRune(word, runeIndex)
	if offset >= 0 {
		w.WriteString(word[:offset])
		w.WriteString(string(highlight))
		w.WriteString(word[offset : offset+width])
		w.WriteString(Reset)
		w.WriteString(string(base))
		w.WriteString(word[offset+width:])
	} else {
		w.WriteString(word)
	}

	w.WriteString(Reset)
}

// returns byte index of a rune index within a string
//
// this is an optimization to allow finding a position of a rune in a string
// and then accessing it directly using string slices, preventing an alloc
func byteIndexOfRune(s string, runeIndex int) (offset, width int) {
	rIdx := 0
	for bIdx, r := range s {
		if rIdx == runeIndex {
			return bIdx, utf8.RuneLen(r)
		}
		rIdx++
	}
	return -1, -1
}

// StrLen ignoring any color codes. If color disabled, returns len(s)
func StrLen(s string) (ret int) {
	if !Enabled {
		return utf8.RuneCountInString(s)
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
