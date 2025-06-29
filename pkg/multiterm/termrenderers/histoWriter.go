package termrenderers

import (
	"fmt"
	"io"
	"strings"

	"github.com/zix99/rare/pkg/color"
	"github.com/zix99/rare/pkg/multiterm"
	"github.com/zix99/rare/pkg/multiterm/termformat"
	"github.com/zix99/rare/pkg/multiterm/termscaler"
	"github.com/zix99/rare/pkg/multiterm/termunicode"
)

type histoPair struct {
	key string
	val int64
}

type HistoWriter struct {
	writer      multiterm.MultilineTerm
	maxVal      int64
	total       int64
	textSpacing int
	items       []histoPair

	ShowBar        bool
	ShowPercentage bool
	Scaler         termscaler.Scaler
	Formatter      termformat.Formatter
}

func NewHistogram(term multiterm.MultilineTerm, maxLines int) *HistoWriter {
	return &HistoWriter{
		writer:         term,
		ShowBar:        true,
		ShowPercentage: true,
		Scaler:         termscaler.ScalerLinear,
		Formatter:      termformat.Default,
		textSpacing:    8,
		items:          make([]histoPair, maxLines),
	}
}

func (s *HistoWriter) WriteFooter(idx int, line string) {
	s.writer.WriteForLine(len(s.items)+idx, line)
}

func (s *HistoWriter) Close() {
	s.writer.Close()
}

func (s *HistoWriter) WriteForLine(line int, key string, val int64) {
	if line > len(s.items) {
		return
	}
	needsFullRefresh := false

	if klen := color.StrLen(key); klen > s.textSpacing {
		s.textSpacing = klen
		needsFullRefresh = true
	}
	if val > s.maxVal {
		s.maxVal = val
		needsFullRefresh = true
	}

	s.items[line] = histoPair{
		key: key,
		val: val,
	}

	if needsFullRefresh {
		s.fullRender()
	} else {
		s.writeLine(line, key, val)
	}
}

func (s *HistoWriter) UpdateTotal(total int64) {
	s.total = total
	s.fullRender()
}

func (s *HistoWriter) fullRender() {
	for idx, item := range s.items {
		if item.val > 0 {
			s.writeLine(idx, item.key, item.val)
		}
	}
}

func (s *HistoWriter) writeLine(line int, key string, val int64) {
	var sb strings.Builder
	sb.Grow(128)

	sb.WriteString(color.Wrapf(color.Yellow, "%-[2]*[1]s", key, s.textSpacing))
	sb.WriteString("    ")
	fmt.Fprintf(&sb, "%-10s", s.Formatter(val, 0, s.maxVal))
	if s.ShowPercentage && s.total > 0 {
		percentage := float64(val) / float64(s.total)
		sb.WriteString(" ")
		sb.WriteString(color.Wrapf(color.Cyan, "[%4.1f%%]", percentage*100.0))
	}

	if s.ShowBar && s.maxVal > 0 {
		sb.WriteString(" ")
		color.Write(&sb, color.Blue, func(w io.StringWriter) {
			termunicode.BarWrite(w, s.Scaler.Scale(val, 0, s.maxVal), 50)
		})
	}

	s.writer.WriteForLine(line, sb.String())
}
