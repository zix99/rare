package markdowncli

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFrontmatterParse(t *testing.T) {
	r := strings.NewReader(`---
summary: hi
order: 1
depth: 2
---
real data
and more real data`)

	fm := ExtractFrontmatter(r)
	assert.Equal(t, "hi", fm.Description())
	assert.Equal(t, 1, fm.Order())
	assert.Equal(t, 2, fm.Depth())
}

func TestEmptyFrontmatter(t *testing.T) {
	r := strings.NewReader(`real data
	and new line`)

	fm := ExtractFrontmatter(r)
	assert.Equal(t, "", fm.Description())
	assert.Equal(t, 0, fm.Order())
	assert.Equal(t, 0, fm.Depth())
}
