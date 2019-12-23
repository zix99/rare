package fuzzy

type LevenshteinKey struct {
	val        []rune
	col        []int
	earlyAbort int
}

func NewLevenshteinKey(val string, earlyAbort float32) *LevenshteinKey {
	return &LevenshteinKey{
		val:        []rune(val),
		col:        make([]int, len(val)+1),
		earlyAbort: len(val) - int(float32(len(val))*earlyAbort),
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

func (s *LevenshteinKey) WordDistance(otherStr string) int {
	other := []rune(otherStr)
	if equalsRunes(s.val, other) {
		return 0
	}

	alen := len(s.val)
	blen := len(other)

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
		if s.col[alen]-(alen-x) > len(s.val)/2 {
			return len(otherStr)
		}
	}
	return s.col[alen]
}

func (s *LevenshteinKey) Distance(other string) float32 {
	dist := s.WordDistance(other)
	sum := len(s.val) + len(other)
	return float32(sum-dist*2) / float32(sum)
}
