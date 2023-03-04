package multiterm

import (
	"io"
	"rare/pkg/multiterm/termstate"
)

var AutoTrim = true

const defaultRows, defaultCols = 24, 80

var computedRows, computedCols = 0, 0

func init() {
	if rows, cols, ok := termstate.GetTermRowsCols(); ok {
		computedRows, computedCols = rows, cols
	} else {
		AutoTrim = false
		computedRows, computedCols = defaultRows, defaultCols
	}
}

func TermRows() int {
	return computedRows
}

func TermCols() int {
	return computedCols
}

func WriteLineNoWrap(out io.Writer, s string) {
	if !AutoTrim {
		out.Write([]byte(s))
		return
	}
	runes := []rune(s)

	visibleRunes := 0
	i := 0
	for i < len(runes) && visibleRunes < computedCols {
		if runes[i] == '\x1b' {
			// parse colors
			for runes[i] != 'm' && i < len(runes)-1 {
				i++
			}
		} else {
			visibleRunes++
		}
		i++
	}

	out.Write([]byte(string(runes[:i])))
}
