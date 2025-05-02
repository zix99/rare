package termformat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringPassthru(t *testing.T) {
	var f StringFormatter = PassthruString
	assert.Equal(t, "hello", f("hello"))
}

func TestStringExpr(t *testing.T) {
	sf, err := StringFromExpression("val {.} or {0} but not {a} or {1}")

	assert.NoError(t, err)
	assert.Equal(t, "val bob or bob but not  or ", sf("bob"))
}

func TestStringExprErr(t *testing.T) {
	_, err := StringFromExpression("{unclosed")
	assert.Error(t, err)
}

// BenchmarkStringExpr-8   	38850990	        30.01 ns/op	      16 B/op	       1 allocs/op
func BenchmarkStringExpr(b *testing.B) {
	sf, _ := StringFromExpression("{.}")

	for range b.N {
		sf("bob")
	}
}
