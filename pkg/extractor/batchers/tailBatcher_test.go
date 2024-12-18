package batchers

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBatchFollowFile(t *testing.T) {
	filenames := make(chan string, 1)
	filenames <- "tailBatcher_test.go" // me

	batcher := TailFilesToChan(filenames, 5, 1, false, false, false)

	batch := <-batcher.BatchChan()
	assert.Equal(t, "tailBatcher_test.go", batch.Source)
	assert.Len(t, batch.Batch, 5)
	assert.NotZero(t, batcher.ReadBytes())
}

func TestBatchFollowTailFile(t *testing.T) {
	tmp, err := os.CreateTemp("", "followtest-")
	if err != nil {
		panic(err)
	}
	defer tmp.Close()
	defer os.Remove(tmp.Name())

	// Add test data
	for i := 0; i < 10; i++ {
		tmp.WriteString("abc\n")
	}

	// Now tail the file
	filenames := make(chan string, 1)
	filenames <- tmp.Name()

	batcher := TailFilesToChan(filenames, 1, 1, false, false, true)

	for batcher.ActiveFileCount() == 0 {
		time.Sleep(time.Millisecond) // Semi-hack: Wait for the go-routine reader to start and the source to be drained
	}

	// And write some more data
	const testLines = 5
	for i := 0; i < testLines; i++ {
		tmp.WriteString("abc\n")
	}

	// And finally assert we got what we wanted
	for i := 0; i < testLines; i++ {
		batch, ok := <-batcher.BatchChan()
		assert.True(t, ok)
		if ok {
			assert.Equal(t, tmp.Name(), batch.Source)
			assert.Equal(t, uint64(i+1), batch.BatchStart)
			assert.Len(t, batch.Batch, 1)
		}
	}

	assert.Len(t, batcher.BatchChan(), 0)
}
