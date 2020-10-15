package aggregation

import (
	"rare/pkg/fuzzy"
)

type FuzzyAggregator struct {
	lookup *fuzzy.FuzzyTable
	Histo  *MatchCounter
}

func NewFuzzyAggregator(matchDist float32, maxOffset int) *FuzzyAggregator {
	return &FuzzyAggregator{
		lookup: fuzzy.NewFuzzyTable(matchDist, maxOffset),
		Histo:  NewCounter(),
	}
}

func (s *FuzzyAggregator) Sample(ele string) {
	_, similarStr, _ := s.lookup.GetMatchId(ele)
	s.Histo.SampleValue(similarStr, 1)
}

func (s *FuzzyAggregator) ParseErrors() uint64 {
	return 0
}
