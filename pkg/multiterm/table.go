package multiterm

import (
	"rare/pkg/color"
	"strings"
)

// A formatted column-row based output with color-coding

type TableWriter struct {
	maxCols, maxRows int
	currentRows      int
	term             *TermWriter
	maxElementLen    int
	rows             [][]string

	HighlightRow0 bool
	HighlightCol0 bool
}

func NewTable(maxCols, maxRows int) *TableWriter {
	return &TableWriter{
		maxCols:       maxCols,
		maxRows:       maxRows,
		term:          New(),
		rows:          make([][]string, maxRows),
		maxElementLen: 8,
		HighlightRow0: true,
		HighlightCol0: true,
	}
}

func (s *TableWriter) InnerWriter() MultilineTerm {
	return s.term
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
		if rowNum == 0 && s.HighlightRow0 {
			sb.WriteString(color.Wrap(color.Underline+color.BrightBlue, cols[i]))
		} else if i == 0 && s.HighlightCol0 {
			sb.WriteString(color.Wrap(color.Yellow, cols[i]))
		} else {
			sb.WriteString(cols[i])
		}
		for j := 0; j < s.maxElementLen-len(cols[i]); j++ {
			sb.WriteRune(' ')
		}
		sb.WriteRune(' ')
	}
	s.term.WriteForLine(rowNum, sb.String())
}
