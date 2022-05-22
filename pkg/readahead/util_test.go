package readahead

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDropCR(t *testing.T) {
	assert.Equal(t, []byte("test"), dropCR([]byte("test")))
	assert.Equal(t, []byte("test\n"), dropCR([]byte("test\n")))
	assert.Equal(t, []byte("test"), dropCR([]byte("test\r")))
}

func TestMaxi(t *testing.T) {
	assert.Equal(t, 1, maxi(0, 1))
	assert.Equal(t, 1, maxi(1, 0))
}
