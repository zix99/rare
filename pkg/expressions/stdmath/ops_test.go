package stdmath

import (
	"maps"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderOfOperationExhaustive(t *testing.T) {
	allOps := slices.Collect(maps.Keys(ops))
	flattenedOrder := slices.Concat(orderOfOps...)
	assert.ElementsMatch(t, allOps, flattenedOrder)
}

func TestOpCodes(t *testing.T) {
	assert.Equal(t, opCodeOrder("*", "+"), -1)
	assert.Equal(t, opCodeOrder("-", "+"), 0)
	assert.Equal(t, opCodeOrder("*", "^"), 1)
}

func TestPrefixInOps(t *testing.T) {
	assert.EqualValues(t, "*", *prefixInOps("*"))
	assert.EqualValues(t, "*", *prefixInOps("*b"))
	assert.EqualValues(t, "<=", *prefixInOps("<= b"))
	assert.EqualValues(t, "<", *prefixInOps("< b"))
}
