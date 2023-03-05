package patterns

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPatternSet(t *testing.T) {
	p := NewPatternSet()
	assert.Len(t, p.patterns, 0)
	p.AddPattern("a", "bb")
	assert.Len(t, p.patterns, 1)
	assert.Len(t, p.Patterns(), 1)
	assert.Equal(t, p.Count(), 1)
}

func TestLoadPatternFile(t *testing.T) {
	p := NewPatternSet()

	r := strings.NewReader("anything .*\n# comment\nanother [1-9]*")
	p.LoadPatternFile(r)

	assert.Equal(t, len(p.patterns), 2)
}
