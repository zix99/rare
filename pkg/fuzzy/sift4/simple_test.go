package sift4

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSameString(t *testing.T) {
	assert.Equal(t, 0, DistanceString("abc def", "abc def", 10))
}

func TestMissingSingle(t *testing.T) {
	assert.Equal(t, 1, DistanceString("abc", "adc", 10))
}

func TestComplexString(t *testing.T) {
	assert.Equal(t, 6, DistanceString("this is a long string with a subtle difference", "this is a short string with a subtle difference", 10))
}

func TestCompletelyDifferent(t *testing.T) {
	assert.Equal(t, 7, DistanceString("abcdefg", "1234567", 10))
}

/*
func TestCompletelyDifferentRatio(t *testing.T) {
	assert.Equal(t, float32(0), DistanceStringRatio("abcdefg", "qqqqqqq"))
}

func TestHalfSimilar(t *testing.T) {
	assert.Equal(t, float32(0.5), DistanceStringRatio("abcdef", "qqqdef"))
}

func TestFullSimilar(t *testing.T) {
	assert.Equal(t, float32(1.0), DistanceStringRatio("abc", "abc"))
}*/

func BenchmarkSimilarityHigh(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DistanceString("this is a very long string to test with", "this is a very short string to test with", 100)
	}
}

func BenchmarkSimilarityLow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DistanceString("this is a very long string to test with", "a completely different string with a few similar words", 5)
	}
}
