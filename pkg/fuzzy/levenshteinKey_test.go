package fuzzy

import (
	"fmt"
	"testing"
)

func BenchmarkKeySimilarityHigh(b *testing.B) {
	key := NewLevenshteinKey("this is a very long string to test with", 2.0)
	var val float32
	for i := 0; i < b.N; i++ {
		val = key.Distance("this is a very short string to test with")
	}
	fmt.Println(val)
}

func BenchmarkKeySimilarityLow(b *testing.B) {
	key := NewLevenshteinKey("this is a very long string to test with", 0.5)
	var val float32
	for i := 0; i < b.N; i++ {
		val = key.Distance("a completely different string with a few similar words")
	}
	fmt.Println(val)
}
