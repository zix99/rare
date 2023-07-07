package termscaler

import (
	"math"
	"strings"
)

type ScalerFunc func(val, min, max int64) float64
type Mapper func(float64) float64

type Scaler struct {
	mapVal   Mapper
	unmapVal Mapper
}

var (
	ScalerNull = Scaler{
		func(f float64) float64 { return 0.0 },
		func(f float64) float64 { return 0.0 },
	}
	ScalerLinear = Scaler{
		func(f float64) float64 { return f },
		func(f float64) float64 { return f },
	}
	ScalerLog2 = Scaler{
		func(f float64) float64 {
			if f <= 1.0 {
				return 0.0
			}
			return math.Log2(f)
		},
		func(f float64) float64 { return math.Pow(2.0, f) },
	}
	ScalerLog10 = Scaler{
		func(f float64) float64 {
			if f <= 1.0 {
				return 0.0
			}
			return math.Log10(f)
		},
		func(f float64) float64 { return math.Pow(10.0, f) },
	}
)

func (s Scaler) remapMinMax(min, max int64) (float64, float64) {
	if max <= min {
		max = min + 1
	}
	return math.Floor(s.mapVal(float64(min))), math.Ceil(s.mapVal(float64(max)))
}

// Returns a val, betwen min and max, to a 0-1 float range
func (s Scaler) Scale(val, min, max int64) float64 {
	if max < min {
		return 0.0
	}
	if val < min {
		return 0.0
	}
	if val > max {
		return 1.0
	}
	minf10, maxf10 := s.remapMinMax(min, max)
	if minf10 == maxf10 {
		return 0.0
	}
	return (s.mapVal(float64(val)) - minf10) / (maxf10 - minf10)
}

// Return [0, bucket-1]
func Bucket(buckets int, unitVal float64) int {
	return int(unitVal * float64(buckets-1))
}

// Return [0, bucket-1]
func (s Scaler) Bucket(buckets int, val, min, max int64) int {
	return Bucket(buckets, s.Scale(val, min, max))
}

// Returns scaled bucket values. May return less buckets than requested for small or high-curve ranges (won't return dupes)
func (s Scaler) ScaleKeys(buckets, min, max int64) []int64 {
	minf10, maxf10 := s.remapMinMax(min, max)

	ret := make([]int64, 0, buckets)
	for i := int64(0); i < buckets; i++ {
		val := int64(s.unmapVal((maxf10-minf10)*float64(i)/float64(buckets-1) + minf10))
		if i == 0 || ret[len(ret)-1] != val {
			ret = append(ret, val)
		}
	}
	return ret
}

func ScalerByName(name string) (Scaler, bool) {
	switch strings.ToLower(name) {
	case "linear", "lin", "":
		return ScalerLinear, true
	case "log10", "log":
		return ScalerLog10, true
	case "log2":
		return ScalerLog2, true
	}
	return ScalerNull, false
}
