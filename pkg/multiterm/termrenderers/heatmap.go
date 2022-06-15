package termrenderers

import (
	"rare/pkg/aggregation"
	"rare/pkg/humanize"
	"rare/pkg/multiterm"
	"rare/pkg/multiterm/termunicode"
	"strings"
)

type Heatmap struct {
	term               multiterm.MultilineTerm
	rowCount, colCount int
	minVal, maxVal     int64
	maxRowKeyWidth     int // Max row width
	currentRows        int // Currently used row count for non-footer
}

func NewHeatmap(term multiterm.MultilineTerm, rows, cols int) *Heatmap {
	return &Heatmap{
		rowCount:       rows,
		colCount:       cols,
		term:           term,
		maxRowKeyWidth: 8,
		maxVal:         1,
	}
}

func (s *Heatmap) WriteTable(agg *aggregation.TableAggregator) {
	// TODO: Correct update of heatmap range
	s.SetMinMax(s.minVal, s.maxVal)

	// Write header
	colNames := agg.OrderedColumnsByName() // TODO: Smart? eg. by number?
	colCount := mini(len(colNames), s.colCount)

	{ // TODO: Make func?
		var sb strings.Builder
		sb.WriteString(strings.Repeat(" ", s.maxRowKeyWidth+1))

		var lastWrite int
		for i := 0; i < colCount; {
			name := colNames[i]
			sb.WriteString(name)
			sb.WriteRune(' ')
			lastWrite = i
			i += len(name) + 1
		}
		if lastWrite < colCount-1 { // TODO: Some smarter context for the last key
			sb.WriteString(colNames[len(colNames)-1])
		}
		s.term.WriteForLine(1, sb.String())
	}

	// Each row...
	rows := agg.OrderedRowsByName()
	for i := 0; i < len(rows); i++ {
		row := rows[i]
		if len(row.Name()) > s.maxRowKeyWidth {
			s.maxRowKeyWidth = len(row.Name())
		}

		var sb strings.Builder
		sb.WriteString(row.Name())
		sb.WriteString(strings.Repeat(" ", s.maxRowKeyWidth-len(row.Name())+1))

		for j := 0; j < colCount; j++ {
			// TODO: Interpolation
			val := row.Value(colNames[j])
			if val < s.minVal {
				s.minVal = val
			}
			if val > s.maxVal {
				s.maxVal = val
			}
			termunicode.HeatWriteLinear(&sb, val, s.minVal, s.maxVal)
		}

		s.term.WriteForLine(2+i, sb.String())
	}

	s.currentRows = 2 + len(rows)
}

func (s *Heatmap) WriteFooter(idx int, line string) {
	s.term.WriteForLine(s.currentRows+idx, line)
}

func (s *Heatmap) SetMinMax(min, max int64) {
	s.minVal = min
	s.maxVal = max

	var sb strings.Builder
	sb.WriteString("        ")

	// Min
	termunicode.HeatWriteLinear(&sb, s.minVal, s.minVal, s.maxVal)
	sb.WriteString(" ")
	sb.WriteString(humanize.Hi(s.minVal))

	// mid-val
	sb.WriteString("    ")
	mid := s.minVal + (s.maxVal-s.minVal)/2
	termunicode.HeatWriteLinear(&sb, mid, s.minVal, s.maxVal)
	sb.WriteString(" ")
	sb.WriteString(humanize.Hi(mid))

	// Max
	sb.WriteString("    ")
	termunicode.HeatWriteLinear(&sb, s.maxVal, s.minVal, s.maxVal)
	sb.WriteString(" ")
	sb.WriteString(humanize.Hi(s.maxVal))

	s.term.WriteForLine(0, sb.String())
}

func mini(i, j int) int {
	if i < j {
		return i
	}
	return j
}
