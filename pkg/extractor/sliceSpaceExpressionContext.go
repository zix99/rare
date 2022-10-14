package extractor

import (
	"rare/pkg/expressions"
	"rare/pkg/expressions/stdlib"
	"rare/pkg/minijson"
	"strconv"
	"strings"
)

type SliceSpaceExpressionContext struct {
	linePtr   string
	indices   []int
	nameTable map[string]int
	source    string
	lineNum   uint64
}

func (s *SliceSpaceExpressionContext) GetMatch(idx int) string {
	sliceIndex := idx * 2
	if sliceIndex < 0 || sliceIndex+1 >= len(s.indices) {
		return ""
	}
	start := s.indices[sliceIndex]
	end := s.indices[sliceIndex+1]
	if start < 0 || end < 0 {
		return ""
	}
	return s.linePtr[start:end]
}

func (s *SliceSpaceExpressionContext) GetKey(key string) string {
	switch key {
	case "src":
		return s.source
	case "line":
		return strconv.FormatUint(s.lineNum, 10)
	case ".":
		return s.json(true, false)
	case "#":
		return s.json(false, true)
	case ".#", "#.":
		return s.json(true, true)
	case "$":
		return s.array()
	}

	if idx, ok := s.nameTable[key]; ok {
		return s.GetMatch(idx)
	}

	return stdlib.ErrorArgName
}

func (s *SliceSpaceExpressionContext) json(named, numbered bool) string {
	var jb minijson.JsonObjectBuilder
	jb.OpenEx(len(s.nameTable) * 50)

	if named {
		for name, idx := range s.nameTable {
			jb.WriteInferred(name, s.GetMatch(idx))
		}
	}
	if numbered {
		for i := 0; i < len(s.indices)/2; i++ {
			if val := s.GetMatch(i); val != "" {
				jb.WriteInferred(strconv.Itoa(i), val)
			}
		}
	}

	jb.Close()

	return jb.String()
}

func (s *SliceSpaceExpressionContext) array() string {
	var sb strings.Builder
	for i := 1; i < len(s.indices)/2; i++ {
		val := s.GetMatch(i)

		if i > 1 {
			sb.WriteRune(expressions.ArraySeparator)
		}
		sb.WriteString(val)
	}
	return sb.String()
}
