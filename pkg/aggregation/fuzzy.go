package aggregation

import (
	"rare/pkg/fuzzy"
)

type FuzzyAggregator struct {
	lookup *fuzzy.FuzzyTable
	Histo  *MatchCounter
}

func NewFuzzyAggregator(matchDist float32, maxOffset, minScore int) *FuzzyAggregator {
	return &FuzzyAggregator{
		lookup: fuzzy.NewFuzzyTable(matchDist, maxOffset, minScore),
		Histo:  NewCounter(),
	}
}

func (s *FuzzyAggregator) Sample(ele string) {
	similarStr, _ := s.lookup.GetMatchId(ele)
	s.Histo.SampleValue(similarStr, 1)
}

func (s *FuzzyAggregator) ParseErrors() uint64 {
	return s.Histo.ParseErrors()
}

func (s *FuzzyAggregator) FuzzyTableSize() int {
	return s.lookup.Count()
}
