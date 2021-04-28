package batchers

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOpenFilesToChan(t *testing.T) {
	filenames := make(chan string, 5)
	filenames <- "fileBatcher_test.go" // me!
	close(filenames)

	batches := OpenFilesToChan(filenames, false, 1, 1)

	total := 0
	var lastStart uint64 = 0
	for batch := range batches {
		assert.Greater(t, batch.BatchStart, lastStart)
		lastStart = batch.BatchStart
		total += len(batch.Batch)
		assert.Equal(t, "fileBatcher_test.go", batch.Source)
	}

	assert.NotZero(t, total)
}

func TestBufferingChan(t *testing.T) {
	var wg sync.WaitGroup

	c := make(chan string)
	wg.Add(1)
	go func() {
		for i := 0; i < 100; i++ {
			c <- "hi"
		}
		close(c)
		wg.Done()
	}()

	bc := bufferChan(c, 100)
	wg.Wait()

	assert.Eventually(t, func() bool {
		return len(bc) == 100
	}, 1*time.Second, 10*time.Millisecond)
}
