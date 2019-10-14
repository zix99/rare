package aggregation

import (
	"math"
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
	min     float64
	max     float64
}

func NewNumericalAggregator() *MatchNumerical {
	return &MatchNumerical{
		values: make([]float64, 0),
		min:    math.MaxFloat64,
		max:    -math.MaxFloat64,
	}
}

func (s *MatchNumerical) Sample(val float64) {
	s.samples++
	s.sum += val
	s.values = append(s.values, val)

	if val < s.min {
		s.min = val
	}
	if val > s.max {
		s.max = val
	}
}

func (s *MatchNumerical) Count() uint64 {
	return s.samples
}

func (s *MatchNumerical) Min() float64 {
	return s.min
}

func (s *MatchNumerical) Max() float64 {
	return s.max
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
