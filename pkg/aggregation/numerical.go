package aggregation

import (
	"math"
	"sort"
	"strconv"
)

type StatisticalAnalysis struct {
	orderedValues []float64
}

type NumericalConfig struct {
	Reverse bool
}

type MatchNumerical struct {
	samples     uint64
	sum         float64
	values      []float64
	min         float64
	max         float64
	parseErrors uint64
	config      *NumericalConfig
}

func NewNumericalAggregator(config *NumericalConfig) *MatchNumerical {
	return &MatchNumerical{
		values: make([]float64, 0),
		min:    math.MaxFloat64,
		max:    -math.MaxFloat64,
		config: config,
	}
}

func (s *MatchNumerical) Samplef(val float64) {
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

func (s *MatchNumerical) Sample(element string) {
	val, err := strconv.ParseFloat(element, 64)
	if err != nil {
		s.parseErrors++
	} else {
		s.Samplef(val)
	}
}

func (s *MatchNumerical) ParseErrors() uint64 {
	return s.parseErrors
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

func (s *MatchNumerical) StdDev() float64 {
	if s.samples == 0 {
		return 0.0
	}
	mean := s.Mean()
	diffSum := 0.0
	for _, v := range s.values {
		diffSum += (v - mean) * (v - mean)
	}
	diffMean := diffSum / float64(s.samples)
	return math.Sqrt(diffMean)
}

func (s *MatchNumerical) Analyze() *StatisticalAnalysis {
	if s.config.Reverse {
		sort.Sort(sort.Reverse(sort.Float64Slice(s.values)))
	} else {
		sort.Float64s(s.values)
	}

	out := &StatisticalAnalysis{
		orderedValues: s.values[0:len(s.values)],
	}

	return out
}

func (s *StatisticalAnalysis) Median() float64 {
	if len(s.orderedValues) == 0 {
		return 0.0
	}
	return s.orderedValues[len(s.orderedValues)/2]
}

func (s *StatisticalAnalysis) Mode() float64 {
	if len(s.orderedValues) == 0 {
		return 0.0
	}
	// We can take advantage of the fact that we know the data
	// here is ordered by counting the max recurrences
	maxObserved := 0
	maxValue := 0.0

	currObserved := 0
	currValue := 0.0

	for i := 0; i < len(s.orderedValues); i++ {
		val := s.orderedValues[i]
		if val != currValue {
			currValue = val
			currObserved = 0
		}
		currObserved++

		if currObserved > maxObserved {
			maxValue = currValue
			maxObserved = currObserved
		}
	}
	return maxValue
}

func (s *StatisticalAnalysis) Quantile(p float64) float64 {
	if len(s.orderedValues) == 0 {
		return 0.0
	}
	idx := int(float64(len(s.orderedValues)) * p)
	return s.orderedValues[idx]
}
