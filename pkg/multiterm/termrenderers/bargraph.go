package termrenderers

import (
	"io"
	"rare/pkg/color"
	"rare/pkg/humanize"
	"rare/pkg/multiterm"
	"rare/pkg/multiterm/termunicode"
	"strings"
)

type BarGraph struct {
	writer multiterm.MultilineTerm

	maxKeyLength int
	subKeys      []string
	maxLineVal   int64
	maxRows      int
	prefixLines  int

	BarSize int
	Stacked bool
}

func NewBarGraph(term multiterm.MultilineTerm) *BarGraph {
	return &BarGraph{
		writer:       term,
		maxKeyLength: 4,
		Stacked:      false,
		BarSize:      50,
	}
}

func (s *BarGraph) SetKeys(keyItems ...string) {
	s.subKeys = keyItems

	if len(keyItems) > 1 || (len(keyItems) == 1 && keyItems[0] != "") {
		s.prefixLines = 1

		var sb strings.Builder
		sb.WriteString(strings.Repeat(" ", s.maxKeyLength+2))
		for idx, item := range keyItems {
			sb.WriteString("  ")
			sb.WriteString(color.Wrap(color.GroupColors[idx%len(color.GroupColors)], termunicode.BarString(1, 1, 1)))
			sb.WriteString(" ")
			sb.WriteString(item)
		}
		s.writer.WriteForLine(0, sb.String())
	}
}

// Writes bar graph values, assuming vals map to the keyItems for each index
func (s *BarGraph) WriteBar(idx int, key string, vals ...int64) {
	if len(key) > s.maxKeyLength {
		s.maxKeyLength = len(key)
	}

	if s.Stacked {
		s.writeBarStacked(idx, key, vals...)
	} else {
		s.writeBarGrouped(idx, key, vals...)
	}
}

func (s *BarGraph) writeBarGrouped(idx int, key string, vals ...int64) {
	for _, val := range vals {
		if val > s.maxLineVal {
			s.maxLineVal = val
		}
	}

	var sb strings.Builder
	sb.WriteString(color.Wrapf(color.Yellow, "%-[2]*[1]s", key, s.maxKeyLength))
	sb.WriteString("  ")

	line := s.prefixLines + idx*len(s.subKeys)

	maxRow := line + len(s.subKeys)
	if maxRow > s.maxRows {
		s.maxRows = maxRow
	}

	for i := 0; i < len(vals); i++ {
		if i > 0 {
			sb.WriteString(strings.Repeat(" ", s.maxKeyLength+2))
		}
		color.Write(&sb, color.GroupColors[i%len(color.GroupColors)], func(w io.StringWriter) {
			termunicode.BarWrite(w, vals[i], s.maxLineVal, int64(s.BarSize))
		})
		sb.WriteString(" ")
		sb.WriteString(humanize.Hi(vals[i]))
		s.writer.WriteForLine(line+i, sb.String())

		sb.Reset()
	}
}

func (s *BarGraph) writeBarStacked(idx int, key string, vals ...int64) {
	var total int64
	for _, val := range vals {
		total += val
	}

	if total > s.maxLineVal {
		s.maxLineVal = total
	}

	var sb strings.Builder
	sb.WriteString(color.Wrapf(color.Yellow, "%-[2]*[1]s", key, s.maxKeyLength))
	sb.WriteString("  ")

	line := idx + s.prefixLines
	maxRow := line + 1
	if maxRow > s.maxRows {
		s.maxRows = maxRow
	}

	for i := 0; i < len(vals); i++ {
		color.Write(&sb, color.GroupColors[i%len(color.GroupColors)], func(w io.StringWriter) {
			termunicode.BarWriteFull(w, vals[i], s.maxLineVal, int64(s.BarSize))
		})
	}
	sb.WriteString("  ")
	sb.WriteString(humanize.Hi(total))
	s.writer.WriteForLine(line, sb.String())
}

func (s *BarGraph) WriteFooter(idx int, str string) {
	s.writer.WriteForLine(s.maxRows+idx, str)
}

func (s *BarGraph) Close() {
	s.writer.Close()
}
