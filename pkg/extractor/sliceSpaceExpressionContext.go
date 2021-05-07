package extractor

import (
	"rare/pkg/expressions/stdlib"
	"strconv"
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
	}

	if idx, ok := s.nameTable[key]; ok {
		return s.GetMatch(idx)
	}

	return stdlib.ErrorArgName
}
