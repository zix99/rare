package humanize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestH(t *testing.T) {
	assert.Equal(t, "Hello 1,000", H("Hello %d", 1000))
}

func TestHDisabled(t *testing.T) {
	Enabled = false
	assert.Equal(t, "Hello 1000", H("Hello %d", 1000))
	Enabled = true
}

func TestHi(t *testing.T) {
	assert.Equal(t, "1,500", Hi(1500))
}

func TestHf(t *testing.T) {
	assert.Equal(t, "1,234,567.8912", Hf(1234567.89121111))
}

func TestHfd(t *testing.T) {
	assert.Equal(t, "1,234,567.89", Hfd(1234567.89121111, 2))
}
