package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReduce(t *testing.T) {
	out, eout, err := testCommandCapture(reduceCommand(),
		`-m (\d+) --snapshot -a "test={sumi {.} {0}}" testdata/log.txt`)
	assert.NoError(t, err)
	assert.Empty(t, eout)
	assert.Equal(t, "test: 32\nMatched: 3 / 6\n96 B (0 B/s) \n", out)
}
