package extractor

type SliceSpaceExpressionContext struct {
	linePtr string
	indices []int
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
