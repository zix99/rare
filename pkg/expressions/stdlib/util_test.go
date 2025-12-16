package stdlib

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPartialString(t *testing.T) {
	assert.True(t, isPartialString("m", "months"))
	assert.True(t, isPartialString("mo", "months"))
	assert.True(t, isPartialString("months", "months"))
	assert.False(t, isPartialString("t", "month"))
	assert.False(t, isPartialString("months", "month"))
}

func TestSimpleParseNumeric(t *testing.T) {
	n, ok := simpleParseNumeric("1234567890")
	assert.True(t, ok)
	assert.Equal(t, int64(1234567890), n)

	_, ok = simpleParseNumeric("1234abc")
	assert.False(t, ok)

	_, ok = simpleParseNumeric("-1234")
	assert.False(t, ok)

	_, ok = simpleParseNumeric("12.34")
	assert.False(t, ok)

	_, ok = simpleParseNumeric("")
	assert.False(t, ok)
}

// BenchmarkSimpleParseNumeric-4   	78927297	        15.44 ns/op	       0 B/op	       0 allocs/op
func BenchmarkSimpleParseNumeric(b *testing.B) {
	for i := 0; i < b.N; i++ {
		simpleParseNumeric("1460653945")
	}
}

// BenchmarkParseInt-4   	24940904	        45.30 ns/op	       0 B/op	       0 allocs/op
func BenchmarkParseInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strconv.ParseInt("1460653945", 10, 64)
	}
}
