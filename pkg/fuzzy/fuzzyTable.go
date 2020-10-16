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
	minScore  int64
	searches  int
}

func NewFuzzyTable(matchDist float32, maxOffset, minScore int) *FuzzyTable {
	return &FuzzyTable{
		keys:      make([]fuzzyItem, 0),
		matchDist: matchDist,
		maxOffset: maxOffset,
		minScore:  int64(minScore),
	}
}

func (s *FuzzyTable) GetMatchId(val string) (id int, match string, isNew bool) {
	for i := range s.keys {
		ele := &s.keys[i]
		d := sift4.DistanceStringRatio(ele.original, val, s.maxOffset)
		if d > s.matchDist {
			ele.score += int64(len(s.keys)) * 2
			return i, ele.original, false
		}
		ele.score--
	}

	s.searches++
	if s.searches >= 10 {
		s.Cleanup()
		s.searches = 0
	}

	newItem := fuzzyItem{
		original: val,
	}
	s.keys = append(s.keys, newItem)

	return len(s.keys) - 1, val, true
}

func (s *FuzzyTable) Cleanup() {
	// Sorting puts the most likely match candidate at the top of the search
	sort.Slice(s.keys, func(i, j int) bool {
		return s.keys[i].score > s.keys[j].score
	})

	// Truncate the list after score falls below threshold
	for i := 0; i < len(s.keys); i++ {
		if s.keys[i].score < s.minScore {
			s.keys = s.keys[:i]
			break
		}
	}
}

func (s *FuzzyTable) Count() int {
	return len(s.keys)
}
