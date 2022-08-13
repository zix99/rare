package termrenderers

import (
	"rare/pkg/aggregation"
	"rare/pkg/color"
	"rare/pkg/humanize"
	"rare/pkg/multiterm"
	"rare/pkg/multiterm/termunicode"
	"strings"
)

type Heatmap struct {
	term               multiterm.MultilineTerm
	rowCount, colCount int
	minVal, maxVal     int64
	FixedMin, FixedMax bool
	maxRowKeyWidth     int // Max row width
	currentRows        int // Currently used row count for non-footer
}

func NewHeatmap(term multiterm.MultilineTerm, rows, cols int) *Heatmap {
	return &Heatmap{
		rowCount:       rows,
		colCount:       cols,
		term:           term,
		maxRowKeyWidth: 0,
		maxVal:         1,
	}
}

func (s *Heatmap) WriteTable(agg *aggregation.TableAggregator) {
	s.UpdateMinMaxFromData(agg)

	// Write header
	colNames := agg.OrderedColumnsByName() // TODO: Smart? eg. by number?
	colCount := s.WriteHeader(colNames)

	// Each row...
	rows := agg.OrderedRowsByName()
	rowCount := mini(len(rows), s.rowCount)
	for i := 0; i < rowCount; i++ {
		s.WriteRow(i, rows[i], colNames[:colCount])
	}

	// If more rows than can display, write how many were missed
	if len(rows) > rowCount {
		s.term.WriteForLine(2+rowCount, color.Wrapf(color.BrightBlack, "(%d more)", len(rows)-rowCount))
		s.currentRows = 3 + rowCount
	} else {
		s.currentRows = 2 + rowCount
	}
}

func (s *Heatmap) WriteFooter(idx int, line string) {
	s.term.WriteForLine(s.currentRows+idx, line)
}

func (s *Heatmap) UpdateMinMaxFromData(agg *aggregation.TableAggregator) {
	min := s.minVal
	if !s.FixedMin {
		min = agg.ComputeMin()
	}
	max := s.maxVal
	if !s.FixedMax {
		max = agg.ComputeMax()
	}

	s.UpdateMinMax(min, max)
}

func (s *Heatmap) UpdateMinMax(min, max int64) {
	s.minVal = min
	s.maxVal = max

	var sb strings.Builder
	for i := 0; i < s.maxRowKeyWidth+1; i++ {
		sb.WriteRune(' ')
	}

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

func (s *Heatmap) WriteHeader(colNames []string) (colCount int) {
	colCount = mini(len(colNames), s.colCount)

	var sb strings.Builder
	sb.WriteString(strings.Repeat(" ", s.maxRowKeyWidth+1))
	const headerDelim = ".."

	for i := 0; i < colCount; {
		name := colNames[i]

		if i != 0 {
			sb.WriteString(headerDelim)
			i += len(headerDelim)
		}

		if i+len(name)+len(headerDelim) > colCount {
			// Too long, jump to last key
			last := colNames[colCount-1]
			indent := colCount - i - len(last)
			if indent > 0 { // Align last name with last col
				sb.WriteString(strings.Repeat(headerDelim[0:1], indent))
				i += indent
			}
			sb.WriteString(underlineHeaderChar(last, colCount-i-1))
			break
		}

		sb.WriteString(underlineHeaderChar(name, 0))
		i += len(name)
	}

	if colCount < len(colNames) {
		sb.WriteString(color.Wrapf(color.BrightBlack, " (%d more)", len(colNames)-s.colCount))
	}

	s.term.WriteForLine(1, sb.String())
	return
}

func (s *Heatmap) WriteRow(idx int, row *aggregation.TableRow, cols []string) {
	if len(row.Name()) > s.maxRowKeyWidth {
		s.maxRowKeyWidth = len(row.Name())
	}

	var sb strings.Builder
	sb.WriteString(color.Wrap(color.Yellow, row.Name()))
	sb.WriteString(strings.Repeat(" ", s.maxRowKeyWidth-len(row.Name())+1))

	for i := 0; i < len(cols); i++ {
		val := row.Value(cols[i])
		termunicode.HeatWriteLinear(&sb, val, s.minVal, s.maxVal)
	}

	s.term.WriteForLine(2+idx, sb.String())
}

func mini(i, j int) int {
	if i < j {
		return i
	}
	return j
}

func underlineHeaderChar(word string, letter int) string {
	if !color.Enabled {
		return word
	}
	if letter >= 0 && letter < len(word) {
		var sb strings.Builder
		sb.Grow(len(word) * 2)
		sb.WriteString(color.Wrap(color.BrightBlue, word[:letter]))
		sb.WriteString(color.Wrap(color.Underline+color.BrightCyan, string(word[letter])))
		sb.WriteString(color.Wrap(color.BrightBlue, word[letter+1:]))
		return sb.String()
	}
	return color.Wrap(color.BrightBlue, word)
}
