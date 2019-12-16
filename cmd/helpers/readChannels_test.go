package helpers

import (
	"rare/pkg/extractor"
	"testing"

	"github.com/hpcloud/tail"
	"github.com/stretchr/testify/assert"
)

func TestTailLineToChan(t *testing.T) {
	tailchan := make(chan *tail.Line)
	ret := tailLineToChan("test", tailchan, 1)

	tailchan <- &tail.Line{
		Text: "Hello",
	}

	val := <-ret
	assert.Equal(t, "test", val.Source)
	assert.Equal(t, extractor.BString("Hello"), val.Batch[0])
	assert.Equal(t, uint64(1), val.BatchStart)

	close(tailchan)
}

func TestOpenFilesToChan(t *testing.T) {
	filenames := make(chan string, 5)
	filenames <- "readChannels_test.go" // me!
	close(filenames)

	batches := openFilesToChan(filenames, false, 1, 1)

	total := 0
	var lastStart uint64 = 0
	for batch := range batches {
		assert.Greater(t, batch.BatchStart, lastStart)
		lastStart = batch.BatchStart
		total += len(batch.Batch)
		assert.Equal(t, "readChannels_test.go", batch.Source)
	}

	assert.NotZero(t, total)
}
