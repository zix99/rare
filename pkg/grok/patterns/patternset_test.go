package patterns

import (
	"bufio"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPatternSet(t *testing.T) {
	p := NewPatternSet()
	assert.Len(t, p.patterns, 0)
	p.AddPattern("a", "bb")
	assert.Len(t, p.patterns, 1)
	assert.Len(t, p.Patterns(), 1)
}

func TestLoadPatternFile(t *testing.T) {
	p := NewPatternSet()

	f, _ := patternFiles.Open("data/common")
	p.LoadPatternFile(bufio.NewReader(f))

	assert.Greater(t, len(p.patterns), 10)
}
