package multiterm

import (
	"fmt"
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
	format      string

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
	needsRefresh := false

	if len(key) > s.textSpacing {
		s.textSpacing = len(key)
		needsRefresh = true
	}
	if val > s.maxVal {
		s.maxVal = val
		needsRefresh = true
	}

	s.items[line] = histoPair{
		key: key,
		val: val,
	}

	if needsRefresh {
		s.format = fmt.Sprintf("%%-%ds    %%-10d  %%s", s.textSpacing)
		for idx, item := range s.items {
			if item.val > 0 {
				progress := item.val * int64(len(progressSlice)) / s.maxVal
				s.writer.WriteForLine(idx, s.format, item.key, item.val, progressSlice[:progress])
			}
		}
	} else {
		progress := val * int64(len(progressSlice)) / s.maxVal
		s.writer.WriteForLine(line, s.format, key, val, progressSlice[:progress])
	}

}
