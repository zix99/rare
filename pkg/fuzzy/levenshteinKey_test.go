package fuzzy

import (
	"log"
	"testing"
)

func BenchmarkKeySimilarityHigh(b *testing.B) {
	key := NewLevenshteinKey("this is a very long string to test with")
	var val float32
	for i := 0; i < b.N; i++ {
		val = key.Distance("this is a very short string to test with", 0)
	}
	log.Println(val)
}

func BenchmarkKeySimilarityLow(b *testing.B) {
	key := NewLevenshteinKey("this is a very long string to test with")
	var val float32
	for i := 0; i < b.N; i++ {
		val = key.Distance("a completely different string with a few similar words", 0.5)
	}
	log.Println(val)
}
