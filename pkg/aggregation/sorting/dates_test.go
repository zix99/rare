package sorting

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDateSort(t *testing.T) {
	vals := []string{"2022-09-03", "2022-09-02", "2021-09-01"}

	Sort(vals, ByDate(func(a, b string) bool {
		panic("fail")
	}))

	assert.Equal(t, []string{"2021-09-01", "2022-09-02", "2022-09-03"}, vals)
}

func TestDateFallback(t *testing.T) {
	vals := []string{"2022-09-03", "2022-09-02", "notadate", "2021-09-01"}

	fellback := false
	Sort(vals, ByDate(func(a, b string) bool {
		fellback = true
		return ByName(a, b)
	}))

	assert.True(t, fellback)
	assert.Equal(t, []string{"2021-09-01", "2022-09-02", "2022-09-03", "notadate"}, vals)
}

func TestByDateWithContextual(t *testing.T) {
	vals := []string{"2022-09-03", "2022-09-02", "2021-09-01"}

	Sort(vals, ByDateWithContextual())

	assert.Equal(t, []string{"2021-09-01", "2022-09-02", "2022-09-03"}, vals)
}
