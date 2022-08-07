package termrenderers

import (
	"rare/pkg/color"
	"rare/pkg/multiterm"
	"strings"
)

// A formatted column-row based output with color-coding

type TableWriter struct {
	maxCols, maxRows int
	activeRows       int
	term             multiterm.MultilineTerm
	colWidth         []int
	rows             [][]string
}

func NewTable(term multiterm.MultilineTerm, maxCols, maxRows int) *TableWriter {
	return &TableWriter{
		maxCols:  maxCols,
		maxRows:  maxRows,
		term:     term,
		rows:     make([][]string, maxRows),
		colWidth: make([]int, maxCols),
	}
}

func (s *TableWriter) WriteFooter(idx int, line string) {
	s.term.WriteForLine(s.activeRows+idx, line)
}

func (s *TableWriter) Close() {
	s.term.Close()
}

func (s *TableWriter) MaxRows() int {
	return s.maxRows
}

func (s *TableWriter) MaxCols() int {
	return s.maxCols
}

func (s *TableWriter) WriteRow(rowNum int, cols ...string) {
	if rowNum >= s.maxRows {
		return
	}
	if rowNum >= s.activeRows {
		s.activeRows = rowNum + 1
	}

	s.rows[rowNum] = cols

	needFullUpdate := false
	for i := 0; i < len(cols) && i < s.maxCols; i++ {
		runeLen := color.StrLen(cols[i])
		if runeLen > s.colWidth[i] {
			s.colWidth[i] = runeLen
			needFullUpdate = true
		}
	}

	if needFullUpdate {
		for i := 0; i < s.activeRows; i++ {
			s.writeRow(i, s.rows[i]...)
		}
	} else {
		s.writeRow(rowNum, cols...)
	}
}

func (s *TableWriter) writeRow(rowNum int, cols ...string) {
	var sb strings.Builder

	for i := 0; i < len(cols) && i < s.maxCols; i++ {
		runeLen := color.StrLen(cols[i])
		sb.WriteString(cols[i])
		for j := 0; j < s.colWidth[i]-runeLen; j++ {
			sb.WriteRune(' ')
		}
		sb.WriteRune(' ')
	}

	s.term.WriteForLine(rowNum, sb.String())
}
