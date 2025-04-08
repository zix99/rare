package stdmath

import (
	"maps"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderOfOperationExhaustive(t *testing.T) {
	allOps := slices.Collect(maps.Keys(ops))
	assert.ElementsMatch(t, allOps, orderOfOps)
}

func TestIsOpBefore(t *testing.T) {
	assert.True(t, isOpAtOrBefore("*", "+"))
	assert.False(t, isOpAtOrBefore("*", "/"))
	assert.False(t, isOpAtOrBefore("*", "*"))
}

func TestOpCodes(t *testing.T) {
	assert.Equal(t, opCodeOrder("*", "+"), 1)
	assert.Equal(t, opCodeOrder("-", "+"), 0)
	assert.Equal(t, opCodeOrder("*", "^"), -1)
}
