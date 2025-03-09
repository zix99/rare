package termformat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpressionFormatter(t *testing.T) {
	f, err := FromExpression("{percent {0} 0}; val: {0}/{val}/{value} min: {1} or {min} max: {2} or {max} blank:{3}{undef}")
	assert.Nil(t, err)
	assert.Equal(t, "500%; val: 5/5/5 min: 0 or 0 max: 10 or 10 blank:", f(5, 0, 10))
}

func TestExpressionFormatterError(t *testing.T) {
	f, err := FromExpression("{unclosed is easy to test")
	assert.Error(t, err)
	assert.Nil(t, f)
}

func TestExpressionMust(t *testing.T) {
	MustFromExpression("{0}")
	assert.Panics(t, func() {
		MustFromExpression("{0")
	})
}

func TestExpressionExpansion(t *testing.T) {
	f, err := FromExpression("bytesize")
	assert.Nil(t, err)
	assert.Equal(t, f(1024, 0, 0), "1 KB")

	f, err = FromExpression("not-a-func")
	assert.Nil(t, err)
	assert.Equal(t, "not-a-func", f(0, 1, 2))
}

// BenchmarkExpression-4   	 6994440	       158.7 ns/op	       8 B/op	       1 allocs/op
func BenchmarkExpression(b *testing.B) {
	f, _ := FromExpression("{0} {1}-{2}")
	for range b.N {
		f(0, 1, 2)
	}
}
