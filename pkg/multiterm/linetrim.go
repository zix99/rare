package multiterm

import (
	"io"
	"os"

	"golang.org/x/term"
)

var AutoTrim = true

const defaultRows, defaultCols = 24, 80

var computedRows, computedCols = 0, 0

func getTermRowsCols() (rows, cols int, ok bool) {
	fd := int(os.Stdout.Fd())
	if !term.IsTerminal(fd) {
		return 0, 0, false
	}

	cols, rows, err := term.GetSize(fd)
	if err != nil {
		return 0, 0, false
	}

	return rows, cols, true
}

func init() {
	if rows, cols, ok := getTermRowsCols(); ok {
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
