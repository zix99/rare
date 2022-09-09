package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildSorter(t *testing.T) {
	assert.NotNil(t, BuildSorterOrFail("text"))
	assert.NotNil(t, BuildSorterOrFail("smart"))
	assert.NotNil(t, BuildSorterOrFail("contextual"))
	assert.NotNil(t, BuildSorterOrFail("value"))
	assert.NotNil(t, BuildSorterOrFail("value:reverse"))
}

func TestDefaultSortResolves(t *testing.T) {
	sortName, _, err := parseSort(DefaultSortFlag.Value)
	assert.NoError(t, err)

	sorter, sorterErr := lookupSorter(sortName)
	assert.NoError(t, sorterErr)
	assert.NotNil(t, sorter)
}
