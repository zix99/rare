package batchers

import (
	"testing"

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
	assert.Contains(t, s.StatusString(), "1/5")
}
