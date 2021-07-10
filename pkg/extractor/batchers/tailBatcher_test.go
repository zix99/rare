package batchers

import (
	"rare/pkg/extractor"
	"testing"

	"github.com/hpcloud/tail"
	"github.com/stretchr/testify/assert"
)

func TestTailLineToChan(t *testing.T) {
	tailchan := make(chan *tail.Line)
	batcher := newBatcher(10)
	go batcher.tailLineToChan("test", tailchan, 1)

	tailchan <- &tail.Line{
		Text: "Hello",
	}

	val := <-batcher.BatchChan()
	assert.Equal(t, "test", val.Source)
	assert.Equal(t, extractor.BString("Hello"), val.Batch[0])
	assert.Equal(t, uint64(1), val.BatchStart)

	close(tailchan)
}

func TestBatchTailFile(t *testing.T) {
	filenames := make(chan string, 1)
	filenames <- "tailBatcher_test.go" // me

	batcher := TailFilesToChan(filenames, 5, false, false)

	batch := <-batcher.c
	assert.Equal(t, "tailBatcher_test.go", batch.Source)
	assert.Len(t, batch.Batch, 5)
	assert.NotZero(t, batcher.ReadBytes())
}
