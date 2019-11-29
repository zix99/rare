package multiterm

import (
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var AutoTrim = true

const defaultRows, defaultCols = 24, 80

func getTermRowsCols() (int, int) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return defaultRows, defaultCols
	}

	parts := strings.Fields(string(out))

	if len(parts) != 2 {
		return defaultRows, defaultCols
	}

	rows, rowsErr := strconv.Atoi(parts[0])
	cols, colsErr := strconv.Atoi(parts[1])
	if rowsErr != nil || colsErr != nil || rows <= 0 || cols <= 0 {
		return defaultRows, defaultCols
	}

	return rows, cols
}

var computedRows, computedCols = 0, 0

func init() {
	if _, ok := os.LookupEnv("TERM"); ok {
		computedRows, computedCols = getTermRowsCols()
	} else {
		AutoTrim = false
		computedRows = defaultRows
		computedCols = defaultCols
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
