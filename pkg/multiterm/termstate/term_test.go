package termstate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPipedOutput(t *testing.T) {
	assert.True(t, IsPipedOutput())
}

func TestGetTerminalSize(t *testing.T) {
	// Can't really test its output.. no idea where it'll be run in a test
	GetTermRowsCols()
}
