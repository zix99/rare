package aggregation

import (
	"sort"
)

type StatisticalAnalysis struct {
	Mean   float64
	Median float64
	Mode   float64
	StdDev float64
}

type MatchNumerical struct {
	samples uint64
	sum     float64
	values  []float64
}

func NewNumericalAggregator() *MatchNumerical {
	return &MatchNumerical{
		values: make([]float64, 0),
	}
}

func (s *MatchNumerical) Sample(val float64) {
	s.samples++
	s.sum += val
	s.values = append(s.values, val)
}

func (s *MatchNumerical) Mean() float64 {
	return s.sum / float64(s.samples)
}

func (s *MatchNumerical) Analyze() *StatisticalAnalysis {
	sort.Float64s(s.values)

	out := &StatisticalAnalysis{}
	if s.samples > 0 {
		out.Mean = s.Mean()
		out.Median = s.values[len(s.values)/2]
		out.Mode = -1

	}

	return out
}
