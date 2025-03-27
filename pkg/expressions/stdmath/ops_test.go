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
