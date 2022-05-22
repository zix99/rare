package testutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func zero(buf []byte) {
	for i := 0; i < len(buf); i++ {
		buf[i] = 0
	}
}

func TestGeneratesData(t *testing.T) {
	buf := make([]byte, 1000)
	rg := NewTextGenerator(50)

	zero(buf)
	n, err := rg.Read(buf[:100])
	assert.Equal(t, n, 50)
	assert.NoError(t, err)
	for i := 0; i < n; i++ {
		assert.NotZero(t, buf[i])
	}

	zero(buf)
	n, err = rg.Read(buf[:10])
	assert.Equal(t, n, 10)
	assert.NoError(t, err)
	for i := 0; i < n; i++ {
		assert.NotZero(t, buf[i])
	}
}
