package levenshtein

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSameString(t *testing.T) {
	assert.Equal(t, 0, DistanceString("abc def", "abc def"))
}

func TestMissingSingle(t *testing.T) {
	assert.Equal(t, 1, DistanceString("abc", "adc"))
}

func TestComplexString(t *testing.T) {
	assert.Equal(t, 4, DistanceString("this is a long string with a subtle difference", "this is a short string with a subtle difference"))
}

func TestCompletelyDifferent(t *testing.T) {
	assert.Equal(t, 7, DistanceString("abcdefg", "1234567"))
}

func TestCompletelyDifferentRatio(t *testing.T) {
	assert.Equal(t, float32(0), DistanceStringRatio("abcdefg", "qqqqqqq"))
}

func TestHalfSimilar(t *testing.T) {
	assert.Equal(t, float32(0.5), DistanceStringRatio("abcdef", "qqqdef"))
}

func TestFullSimilar(t *testing.T) {
	assert.Equal(t, float32(1.0), DistanceStringRatio("abc", "abc"))
}

func BenchmarkSimilarity(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DistanceStringRatio("abcdef", "qqqdef")
	}
}
