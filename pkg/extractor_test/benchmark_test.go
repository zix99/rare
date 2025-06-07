package benchmark_test

import (
	"testing"

	"github.com/zix99/rare/pkg/extractor"
	"github.com/zix99/rare/pkg/matchers"
	"github.com/zix99/rare/pkg/matchers/fastregex"
)

func batchInputGenerator(batches int, batchSize int) <-chan extractor.InputBatch {
	c := make(chan extractor.InputBatch, 128)
	go func() {
		for i := 0; i < batches; i++ {
			batch := make([]extractor.BString, batchSize)
			for j := 0; j < batchSize; j++ {
				batch[j] = extractor.BString("abcdefg 123")
			}
			c <- extractor.InputBatch{
				Batch:      batch,
				Source:     "",
				BatchStart: 1,
			}
		}
		close(c)
	}()
	return c
}

func BenchmarkExtractor(b *testing.B) {
	total := 0
	for n := 0; n < b.N; n++ {
		gen := batchInputGenerator(10000, 100)
		extractor, _ := extractor.New(gen, &extractor.Config{
			Matcher: matchers.ToFactory(fastregex.MustCompile(`(\d{3})`)),
			Extract: "{bucket {1} 10}",
			Workers: 2,
		})
		reader := extractor.ReadFull()
		for val := range reader {
			total++
			if val[0].Extracted != "120" {
				panic("NO MATCH")
			}
		} // Drain reader
	}
}
