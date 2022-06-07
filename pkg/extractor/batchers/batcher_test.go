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

	s.stopFileReading("abc")
	assert.Equal(t, 5, s.sourceCount)
	assert.Equal(t, 1, s.readCount)
	assert.Equal(t, 0, s.ReadErrors())
	assert.Contains(t, s.StatusString(), "1/5")
}

func TestReaderToBatcher(t *testing.T) {
	s := newBatcher(10)

	testData := `line1
line2
line3`

	s.syncReaderToBatcher("string", strings.NewReader(testData), 2)

	b1 := <-s.BatchChan()
	b2 := <-s.BatchChan()

	assert.Len(t, b1.Batch, 2)
	assert.Len(t, b2.Batch, 1)
	assert.Equal(t, s.errorCount, 0)
	assert.Equal(t, s.ReadBytes(), uint64(17))
}

func TestBatcherWithAutoFlush(t *testing.T) {
	s := newBatcher(10)

	testData := `line1
line2
line3`

	s.syncReaderToBatcherWithTimeFlush("string", strings.NewReader(testData), 2, 1*time.Second)

	b1 := <-s.BatchChan()
	b2 := <-s.BatchChan()

	assert.Len(t, b1.Batch, 2)
	assert.Len(t, b2.Batch, 1)
	assert.Equal(t, s.errorCount, 0)
	assert.Equal(t, s.ReadBytes(), uint64(17))
}
