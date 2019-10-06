package multiterm

import (
	"strings"
)

type histoPair struct {
	key string
	val int64
}

type HistoWriter struct {
	writer      *TermWriter
	maxVal      int64
	textSpacing int
	items       []histoPair

	ShowBar bool
}

func NewHistogram(maxLines int) *HistoWriter {
	return &HistoWriter{
		writer:      New(maxLines),
		ShowBar:     true,
		textSpacing: 16,
		items:       make([]histoPair, maxLines),
	}
}

var progressSlice string = strings.Repeat("|", 50)

func (s *HistoWriter) WriteForLine(line int, key string, val int64) {
	if line > len(s.items) {
		return
	}
	needsFullRefresh := false

	if len(key) > s.textSpacing {
		s.textSpacing = len(key)
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
		for idx, item := range s.items {
			if item.val > 0 {
				s.writeLine(idx, item.key, item.val)
			}
		}
	} else {
		s.writeLine(line, key, val)
	}
}

func (s *HistoWriter) writeLine(line int, key string, val int64) {
	if s.ShowBar {
		progress := val * int64(len(progressSlice)) / s.maxVal
		s.writer.WriteForLine(line, "%-[4]*[1]s    %-10[2]d %[3]s", key, val, progressSlice[:progress], s.textSpacing)
	} else {
		s.writer.WriteForLine(line, "%-[3]*[1]s    %-10[2]d", key, val, s.textSpacing)
	}
}
