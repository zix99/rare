package sorting

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNameValueSorter(t *testing.T) {
	arr := []NameValuePair{
		{"b", 123},
		{"q", 44},
		{"a", 44},
	}

	Sort(arr, ValueSorter)

	expected := []NameValuePair{
		{"b", 123},
		{"a", 44},
		{"q", 44},
	}
	assert.Equal(t, expected, arr)
}

func TestNameValueNilSorter(t *testing.T) {
	arr := []NameValuePair{
		{"b", 123},
		{"q", 44},
		{"a", 44},
	}

	Sort(arr, ValueNilSorter(ByName))

	expected := []NameValuePair{
		{"a", 44},
		{"b", 123},
		{"q", 44},
	}
	assert.Equal(t, expected, arr)
}
