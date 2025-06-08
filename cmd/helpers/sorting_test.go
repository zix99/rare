package helpers

import (
	"testing"

	"github.com/zix99/rare/pkg/aggregation/sorting"

	"github.com/stretchr/testify/assert"
)

func TestBuildSorter(t *testing.T) {
	assert.NotNil(t, BuildSorterOrFail("text"))
	assert.NotNil(t, BuildSorterOrFail("numeric"))
	assert.NotNil(t, BuildSorterOrFail("contextual"))
	assert.NotNil(t, BuildSorterOrFail("value"))
	assert.NotNil(t, BuildSorterOrFail("value:reverse"))
	testLogFatal(t, 2, func() {
		BuildSorterOrFail("fake")
	})
}

func TestOrderResults(t *testing.T) {
	assertSortEquals(t, "text", 1, 4, 2, 0, 3)
	assertSortEquals(t, "text:asc", 1, 4, 2, 0, 3)
	assertSortEquals(t, "text:reverse", 3, 0, 2, 4, 1)
	assertSortEquals(t, "text:desc", 3, 0, 2, 4, 1)

	assertSortEquals(t, "numeric", 1, 4, 2, 0, 3)
	assertSortEquals(t, "numeric:asc", 1, 4, 2, 0, 3)
	assertSortEquals(t, "numeric:reverse", 3, 0, 2, 4, 1)
	assertSortEquals(t, "numeric:desc", 3, 0, 2, 4, 1)

	assertSortEquals(t, "value", 3, 2, 1, 0, 4)
	assertSortEquals(t, "value:desc", 3, 2, 1, 0, 4)
	assertSortEquals(t, "value:reverse", 4, 0, 1, 2, 3)
	assertSortEquals(t, "value:asc", 4, 0, 1, 2, 3)
}

func TestInvalidSortNames(t *testing.T) {
	sorter, err := BuildSorter("bla")
	assert.Nil(t, sorter)
	assert.Error(t, err)

	sorter, err = BuildSorter("numeric:bla")
	assert.Nil(t, sorter)
	assert.Error(t, err)
}

// Given a hardcoded set of values, and a sort name assert the order is as expected
func assertSortEquals(t *testing.T, sortName string, order ...int) {
	sorter, err := BuildSorter(sortName)
	assert.NoError(t, err)

	type orderedPair struct {
		sorting.NameValuePair
		id int
	}

	vals := []orderedPair{
		{sorting.NameValuePair{Name: "qef", Value: 5}, 0},
		{sorting.NameValuePair{Name: "abc", Value: 12}, 1},
		{sorting.NameValuePair{Name: "egf", Value: 52}, 2},
		{sorting.NameValuePair{Name: "zac", Value: 52}, 3},
		{sorting.NameValuePair{Name: "bbb", Value: 3}, 4},
	}

	if len(order) != len(vals) {
		panic("bad test")
	}

	sorting.SortBy(vals, sorter, func(obj orderedPair) sorting.NameValuePair {
		return obj.NameValuePair
	})

	for i := 0; i < len(vals); i++ {
		assert.Equal(t, order[i], vals[i].id)
	}

}

func TestDefaultSortResolves(t *testing.T) {
	sortName, _, err := parseSort(DefaultSortFlag.Value)
	assert.NoError(t, err)

	sorter, sorterErr := lookupSorter(sortName)
	assert.NoError(t, sorterErr)
	assert.NotNil(t, sorter)
}

func TestBuildSortFlag(t *testing.T) {
	flag := DefaultSortFlagWithDefault("contextual")
	assert.Equal(t, "contextual", flag.Value)

	assert.Panics(t, func() {
		DefaultSortFlagWithDefault("fake")
	})
}
