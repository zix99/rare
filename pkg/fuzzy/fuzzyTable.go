package fuzzy

import (
	"rare/pkg/fuzzy/sift4"
	"sort"
)

type fuzzyItem struct {
	original string
	score    int64
}

type FuzzyTable struct {
	keys      []fuzzyItem
	matchDist float32
	maxOffset int
	maxSize   int
	searches  int
}

func NewFuzzyTable(matchDist float32, maxOffset, maxSize int) *FuzzyTable {
	if maxSize < 0 {
		panic("Invalid max size")
	}
	if maxOffset < 0 {
		panic("Invalid max offset")
	}
	return &FuzzyTable{
		keys:      make([]fuzzyItem, 0),
		matchDist: matchDist,
		maxOffset: maxOffset,
		maxSize:   maxSize,
	}
}

func (s *FuzzyTable) GetMatchId(val string) (match string, isNew bool) {
	for i := range s.keys {
		ele := &s.keys[i]
		d := sift4.DistanceStringRatio(ele.original, val, s.maxOffset)
		if d > s.matchDist {
			if d < 0.99 { // Imperfect matches score more
				ele.score += int64(len(s.keys))
			} else {
				ele.score++
			}
			return ele.original, false
		}
		ele.score--
	}

	s.searches++
	if s.searches >= 10 {
		s.Cleanup()
		s.searches = 0
	}

	if len(s.keys) < s.maxSize || s.keys[len(s.keys)-1].score < 1 {
		newItem := fuzzyItem{
			original: val,
			score:    1,
		}
		s.keys = append(s.keys, newItem)
	}

	return val, true
}

func (s *FuzzyTable) Cleanup() {
	// Sorting puts the most likely match candidate at the top of the search
	sort.Slice(s.keys, func(i, j int) bool {
		return s.keys[i].score > s.keys[j].score
	})

	if len(s.keys) > s.maxSize {
		s.keys = s.keys[:s.maxSize]
	}
}

func (s *FuzzyTable) Count() int {
	return len(s.keys)
}
