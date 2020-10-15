package fuzzy

import (
	"rare/pkg/fuzzy/sift4"
)

type fuzzyItem struct {
	original string
}

type FuzzyTable struct {
	keys      []fuzzyItem
	matchDist float32
	maxOffset int
}

func NewFuzzyTable(matchDist float32, maxOffset int) *FuzzyTable {
	return &FuzzyTable{
		keys:      make([]fuzzyItem, 0),
		matchDist: matchDist,
		maxOffset: maxOffset,
	}
}

func (s *FuzzyTable) GetMatchId(val string) (id int, match string, isNew bool) {
	for idx, ele := range s.keys {
		d := sift4.DistanceStringRatio(ele.original, val, s.maxOffset)
		if d > s.matchDist {
			return idx, ele.original, false
		}
	}

	newItem := fuzzyItem{
		original: val,
	}
	s.keys = append(s.keys, newItem)

	return len(s.keys) - 1, val, true
}

func (s *FuzzyTable) Count() int {
	return len(s.keys)
}
