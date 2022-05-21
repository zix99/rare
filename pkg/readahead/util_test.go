package readahead

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDropCR(t *testing.T) {
	r := strings.NewReader("test\r\nthing")
	ra := NewBuffered(r, 3)
	assert.Equal(t, []byte("test"), ra.ReadLine())
	assert.Equal(t, []byte("thing"), ra.ReadLine())
	assert.Nil(t, ra.ReadLine())
}

func TestMaxi(t *testing.T) {
	assert.Equal(t, 1, maxi(0, 1))
	assert.Equal(t, 1, maxi(1, 0))
}
