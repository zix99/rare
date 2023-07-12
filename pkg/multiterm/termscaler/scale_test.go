package termscaler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinearScale(t *testing.T) {
	s, _ := ScalerByName("linear")
	assert.Equal(t, 0, s.Bucket(5, -5, 0, 10))
	assert.Equal(t, 0, s.Bucket(5, 0, 0, 10))
	assert.Equal(t, 2, s.Bucket(5, 5, 0, 10))
	assert.Equal(t, 4, s.Bucket(5, 10, 0, 10))
	assert.Equal(t, 4, s.Bucket(5, 20, 0, 10))
	assert.Equal(t, 0, s.Bucket(5, 120, 100, 200))
	assert.Equal(t, 3, s.Bucket(5, 175, 100, 200))

	// Edge cases
	assert.Equal(t, 0, s.Bucket(5, 20, 20, 10))
	assert.Equal(t, 4, s.Bucket(5, 20, 10, 10))
	assert.Equal(t, 0, s.Bucket(5, -100, 0, -10))
}

func TestLogScale(t *testing.T) {
	s, _ := ScalerByName("log10")
	mmin, mmax := s.remapMinMax(1, 10000)
	assert.Equal(t, 0.0, mmin)
	assert.Equal(t, 4.0, mmax)

	assert.Equal(t, 0, s.Bucket(5, -100, 0, 10000))
	assert.Equal(t, 0, s.Bucket(5, 0, 0, 10000))
	assert.Equal(t, 0, s.Bucket(5, 3, 0, 10000))
	assert.Equal(t, 0, s.Bucket(5, 7, 0, 10000))
	assert.Equal(t, 1, s.Bucket(5, 10, 0, 10000))
	assert.Equal(t, 2, s.Bucket(5, 100, 0, 10000))
	assert.Equal(t, 3, s.Bucket(5, 1000, 0, 10000))
	assert.Equal(t, 3, s.Bucket(5, 2000, 0, 10000))
	assert.Equal(t, 3, s.Bucket(5, 7000, 0, 10000))
	assert.Equal(t, 4, s.Bucket(5, 10000, 0, 10000))
	assert.Equal(t, 4, s.Bucket(5, 20000, 0, 10000))

	assert.Equal(t, 0, s.Bucket(5, 10, 100, 10000))
	assert.Equal(t, 0, s.Bucket(5, 100, 100, 10000))
	assert.Equal(t, 2, s.Bucket(5, 1000, 100, 10000))
	assert.Equal(t, 3, s.Bucket(5, 5000, 100, 10000))
	assert.Equal(t, 4, s.Bucket(5, 10000, 100, 10000))

	// Edge cases
	assert.Equal(t, 4, s.Bucket(5, 20000, 100, 100))
	assert.Equal(t, 0, s.Bucket(5, 20000, 100, -100))
	assert.Equal(t, 0, s.Bucket(0, 0, -100, 100))
	assert.Equal(t, 0, s.Bucket(0, 1, -100, 100))
}

// More realistic use-case
func TestLogScale2(t *testing.T) {
	s := ScalerLog10

	assert.Equal(t, 0, s.Bucket(16, -100, 0, 10000))
	assert.Equal(t, 0, s.Bucket(16, 0, 0, 10000))
	assert.Equal(t, 1, s.Bucket(16, 2, 0, 10000))
	assert.Equal(t, 1, s.Bucket(16, 3, 0, 10000))
	assert.Equal(t, 1, s.Bucket(16, 2, 0, 10000))
	assert.Equal(t, 2, s.Bucket(16, 5, 0, 10000))
	assert.Equal(t, 3, s.Bucket(16, 7, 0, 10000))
	assert.Equal(t, 14, s.Bucket(16, 9999, 0, 10000))
	assert.Equal(t, 15, s.Bucket(16, 10000, 0, 10000))
	assert.Equal(t, []int64{1, 3, 6, 11, 21, 39, 73, 135, 251, 464, 857, 1584, 2928, 5411, 10000}, s.ScaleKeys(16, 0, 10000))
}

func TestLinearKeySet(t *testing.T) {
	assert.Equal(t, []int64{0, 25, 50, 75, 100}, ScalerLinear.ScaleKeys(5, 0, 100))
	assert.Equal(t, []int64{50, 62, 75, 87, 100}, ScalerLinear.ScaleKeys(5, 50, 100))
	assert.Equal(t, []int64{0, 1, 2}, ScalerLinear.ScaleKeys(10, 0, 2))
	assert.Equal(t, []int64{-10, -7, -5, -3, -1, 1, 3, 5, 7, 10}, ScalerLinear.ScaleKeys(10, -10, 10))
}

func TestLog10KeySet(t *testing.T) {
	assert.Equal(t, []int64{1, 10, 100, 1000, 10000}, ScalerLog10.ScaleKeys(5, 0, 10000))
	assert.Equal(t, []int64{100, 1000, 10000}, ScalerLog10.ScaleKeys(3, 100, 10000))
	assert.Equal(t, []int64{1, 3, 10, 31, 100}, ScalerLog10.ScaleKeys(5, 0, 100))
	assert.Equal(t, []int64{10, 17, 31, 56, 100}, ScalerLog10.ScaleKeys(5, 50, 100))
	assert.Equal(t, []int64{1, 3, 10, 31, 100}, ScalerLog10.ScaleKeys(5, -100, 100))
}

func TestLog2KeySet(t *testing.T) {
	s, _ := ScalerByName("log2")
	assert.Equal(t, []int64{1, 11, 128, 1448, 16384}, s.ScaleKeys(5, 0, 10000))
	assert.Equal(t, []int64{64, 1024, 16384}, s.ScaleKeys(3, 100, 10000))
	assert.Equal(t, []int64{1, 3, 11, 38, 128}, s.ScaleKeys(5, 0, 100))
	assert.Equal(t, []int64{32, 45, 64, 90, 128}, s.ScaleKeys(5, 50, 100))
	assert.Equal(t, []int64{1, 3, 11, 38, 128}, s.ScaleKeys(5, -100, 100))
}

func TestInvalidScalerByName(t *testing.T) {
	s, ok := ScalerByName("fake-name")
	assert.False(t, ok)
	assert.NotNil(t, s.mapVal)
	assert.NotNil(t, s.unmapVal)
	assert.Equal(t, 0.0, s.Scale(0, 0, 0))
	assert.Equal(t, 0.0, s.Scale(0, 0, 10))
	assert.Equal(t, []int64{0}, s.ScaleKeys(10, 0, 2))
}
