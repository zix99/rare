package sorting

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDayOfWeekSort(t *testing.T) {
	list := []string{
		"wed",
		"tues",
		"mon",
		"thurs",
	}
	sorter := ByContextual()
	Sort(list, sorter)

	assert.Equal(t, []string{"mon", "tues", "wed", "thurs"}, list)
}

func TestFallbackSort(t *testing.T) {
	list := []string{"wed", "abc", "00"}
	sorter := ByContextual()
	Sort(list, sorter)

	assert.Equal(t, []string{"00", "abc", "wed"}, list)
}
