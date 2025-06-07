package helpers

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/zix99/rare/pkg/extractor"
	"github.com/zix99/rare/pkg/extractor/batchers"
	"github.com/zix99/rare/pkg/matchers"
	"github.com/zix99/rare/pkg/matchers/fastregex"

	"github.com/stretchr/testify/assert"
)

var testData = `abc 123
def 245
qqq 123
xxx`

type VirtualAggregator struct {
	items []string
}

func (s *VirtualAggregator) Sample(element string) {
	s.items = append(s.items, element)
}

func (s *VirtualAggregator) ParseErrors() uint64 {
	return 0
}

func TestAggregationLoop(t *testing.T) {
	// Build a real extractor
	batcher := batchers.OpenReaderToChan("test", io.NopCloser(strings.NewReader(testData)), 1, 1)
	ex, err := extractor.New(batcher.BatchChan(), &extractor.Config{
		Matcher: matchers.ToFactory(fastregex.MustCompile(`(\d+)`)),
		Extract: "val:{1}",
		Workers: 1,
	})
	assert.NoError(t, err)

	// Build a fake aggregator
	agg := &VirtualAggregator{}

	outputTriggered := 0
	RunAggregationLoop(ex, agg, func() {
		outputTriggered++
	})

	// Validation
	assert.GreaterOrEqual(t, outputTriggered, 1)
	assert.Equal(t, 3, len(agg.items))

	// Also validate summary building since we have all the correct context
	WriteExtractorSummary(ex)
	FWriteExtractorSummary(ex, 1, "abc")
	WriteBatcherSummary(os.Stderr, batcher, nil)
}
