package batchers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBatchTailFile(t *testing.T) {
	filenames := make(chan string, 1)
	filenames <- "tailBatcher_test.go" // me

	batcher := TailFilesToChan(filenames, 5, false, false)

	batch := <-batcher.c
	assert.Equal(t, "tailBatcher_test.go", batch.Source)
	assert.Len(t, batch.Batch, 5)
	assert.NotZero(t, batcher.ReadBytes())
}
