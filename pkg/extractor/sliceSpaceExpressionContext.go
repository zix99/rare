package extractor

type SliceSpaceExpressionContext struct {
	linePtr string
	indices []int
}

func (s *SliceSpaceExpressionContext) GetMatch(idx int) string {
	start := idx * 2
	if start < 0 || start+1 >= len(s.indices) {
		return ""
	}
	return s.linePtr[s.indices[start]:s.indices[start+1]]
}
