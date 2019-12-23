package fuzzy

type LevenshteinKey struct {
	val []rune
	col []int
}

func NewLevenshteinKey(val string) *LevenshteinKey {
	return &LevenshteinKey{
		val: []rune(val),
		col: make([]int, len(val)+1),
	}
}

func equalsRunes(a, b []rune) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func (s *LevenshteinKey) WordDistance(otherStr string, abortAt float32) int {
	other := []rune(otherStr)
	if equalsRunes(s.val, other) {
		return 0
	}

	alen := len(s.val)
	blen := len(other)

	abortDistance := int(float32(alen) * (1.0 - abortAt))

	for y := 1; y <= alen; y++ {
		s.col[y] = y
	}

	for x := 1; x <= blen; x++ {
		s.col[0] = x
		lastkey := x - 1
		for y := 1; y <= alen; y++ {
			oldkey := s.col[y]
			var incr int
			if s.val[y-1] != other[x-1] {
				incr = 1
			}
			s.col[y] = min3(s.col[y]+1, s.col[y-1]+1, lastkey+incr)
			lastkey = oldkey
		}
		approxDist := s.col[alen] - (alen - x)
		if approxDist*2 > abortDistance {
			return len(otherStr)
		}
	}
	return s.col[alen]
}

func (s *LevenshteinKey) Distance(other string, abortAt float32) float32 {
	dist := s.WordDistance(other, abortAt)
	sum := len(s.val) + len(other)
	return float32(sum-dist*2) / float32(sum)
}
