package multiterm

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTermSizeIsSet(t *testing.T) {
	assert.Greater(t, TermRows(), 1)
	assert.Greater(t, TermCols(), 1)
}

func TestWriteLineWrap(t *testing.T) {
	computedCols = 10
	WriteLineNoWrap(os.Stdout, "hello there this \x1b123m is a longer than 10 char string")
}
