package humanize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHDisabled(t *testing.T) {
	Enabled = false
	assert.Equal(t, "1000", Hi(1000))
	assert.Equal(t, "1000", Hui(1000))
	assert.Equal(t, "1000", Hi32(1000))
	assert.Equal(t, "1000.0000", Hf(1000.0))
	assert.Equal(t, "1000.00000", Hfd(1000.0, 5))
	assert.Equal(t, "12341234", ByteSize(12341234))
	Enabled = true
}

func TestHi(t *testing.T) {
	assert.Equal(t, "1,500", Hi(1500))
}

func TestHui(t *testing.T) {
	assert.Equal(t, "1,500", Hui(1500))
}

func TestHi32(t *testing.T) {
	assert.Equal(t, "1,500", Hi32(1500))
}

func TestHf(t *testing.T) {
	assert.Equal(t, "1,234,567.8912", Hf(1234567.89121111))
}

func TestHfd(t *testing.T) {
	assert.Equal(t, "1,234,567.89", Hfd(1234567.89121111, 2))
}
