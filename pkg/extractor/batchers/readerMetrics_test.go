package batchers

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReaderMetrics(t *testing.T) {
	r := newReaderMetrics(strings.NewReader("abc"))
	buf := make([]byte, 2)
	r.Read(buf)

	assert.Equal(t, uint64(2), r.readBytes)
	assert.Equal(t, uint64(2), r.CountReset())
	assert.Equal(t, uint64(0), r.CountReset())
}
