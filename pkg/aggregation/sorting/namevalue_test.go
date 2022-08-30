package sorting

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type pair struct {
	s string
	v int64
}

func (s pair) SortName() string {
	return s.s
}
func (s pair) SortValue() int64 {
	return s.v
}

func TestNameValueSorter(t *testing.T) {
	arr := []pair{
		{"b", 123},
		{"q", 44},
		{"a", 44},
	}

	SortNameValue(arr, ValueSorter())

	expected := []pair{
		{"b", 123},
		{"a", 44},
		{"q", 44},
	}
	assert.Equal(t, expected, arr)
}

func TestNameValueNilSorter(t *testing.T) {
	arr := []pair{
		{"b", 123},
		{"q", 44},
		{"a", 44},
	}

	SortNameValue(arr, ValueNameSorter(ByName))

	expected := []pair{
		{"a", 44},
		{"b", 123},
		{"q", 44},
	}
	assert.Equal(t, expected, arr)
}
