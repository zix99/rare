package multiterm

import (
	"fmt"
	"strings"
)

// A formatted column-row based output with color-coding

type TableWriter struct {
	maxCols, maxRows int
	currentRows      int
	term             *TermWriter
	maxElementLen    int
	rows             [][]string
}

func NewTable(maxCols, maxRows int) *TableWriter {
	return &TableWriter{
		maxCols:       maxCols,
		maxRows:       maxRows,
		term:          New(maxRows),
		rows:          make([][]string, maxRows),
		maxElementLen: 10,
	}
}

func (s *TableWriter) WriteRow(rowNum int, cols ...string) {
	if rowNum >= s.maxRows {
		return
	}
	if rowNum > s.currentRows {
		s.currentRows = rowNum
	}

	s.rows[rowNum] = cols

	needFullUpdate := false
	for _, val := range cols {
		if len(val) > s.maxElementLen {
			s.maxElementLen = len(val)
		}
	}

	if needFullUpdate {
		for i := 0; i < s.currentRows; i++ {
			s.writeRow(i, s.rows[i]...)
		}
	} else {
		s.writeRow(rowNum, cols...)
	}
}

func (s *TableWriter) writeRow(rowNum int, cols ...string) {
	var sb strings.Builder
	for i := 0; i < len(cols) && i < s.maxCols; i++ {
		sb.WriteString(fmt.Sprintf("%-[1]*[2]s", s.maxElementLen+1, cols[i]))
	}
	s.term.WriteForLine(rowNum, sb.String())
}
