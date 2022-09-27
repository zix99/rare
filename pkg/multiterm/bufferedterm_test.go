package multiterm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPipedOutput(t *testing.T) {
	assert.True(t, IsPipedOutput())
}

func TestBufferedTerm(t *testing.T) {
	vt := NewBufferedTerm()
	vt.WriteForLine(0, "hello")
	vt.Close()
}
