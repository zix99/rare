package batchers

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBatcherTracking(t *testing.T) {
	s := newBatcher(1)

	s.setSourceCount(5)

	assert.Equal(t, 5, s.sourceCount)

	s.startFileReading("abc")
	assert.Contains(t, s.StatusString(), "0/5")
	assert.Equal(t, s.ActiveFileCount(), 1)

	s.stopFileReading("abc")
	assert.Equal(t, 5, s.sourceCount)
	assert.Equal(t, 1, s.ReadFiles())
	assert.Equal(t, 0, s.ReadErrors())
	assert.Contains(t, s.StatusString(), "1/5")
	assert.Equal(t, s.ActiveFileCount(), 0)
}

func TestReaderToBatcher(t *testing.T) {
	s := newBatcher(10)

	testData := `line1
line2
line3`

	s.syncReaderToBatcher("string", strings.NewReader(testData), 2, 1024)

	b1 := <-s.BatchChan()
	b2 := <-s.BatchChan()

	assert.Len(t, b1.Batch, 2)
	assert.Len(t, b2.Batch, 1)
	assert.Equal(t, s.errorCount, 0)
	assert.Equal(t, s.ReadBytes(), uint64(17))
	assert.Equal(t, s.ActiveFileCount(), 0)
}

func TestBatcherWithAutoFlush(t *testing.T) {
	s := newBatcher(10)

	testData := `line1
line2
line3`

	s.syncReaderToBatcherWithTimeFlush("string", strings.NewReader(testData), 2, 1024, 1*time.Second)

	b1 := <-s.BatchChan()
	b2 := <-s.BatchChan()

	assert.Len(t, b1.Batch, 2)
	assert.Len(t, b2.Batch, 1)
	assert.Equal(t, s.errorCount, 0)
	assert.Equal(t, s.ReadBytes(), uint64(17))
}

func TestDurationFormat(t *testing.T) {
	assert.Equal(t, "020ms", durationToString(20*time.Millisecond))
	assert.Equal(t, "1.12s", durationToString(1120*time.Millisecond))
	assert.Equal(t, "35m2.0s", durationToString(35*time.Minute+2*time.Second))
}
