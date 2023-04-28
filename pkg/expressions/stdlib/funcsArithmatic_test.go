package stdlib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArithmatic(t *testing.T) {
	kb, _ := NewStdKeyBuilder().Compile("{sumi {1} {4}} {multi {1} 2} {divi {1} 2} {subi {1} 10}")
	key := kb.BuildKey(mockContext("ab", "100", "1000000", "5000000.123456", "22"))
	assert.Equal(t, "122 200 50 90", key)
}

func TestArithmaticError(t *testing.T) {
	kb, _ := NewStdKeyBuilder().Compile("{sumi 1} {sumi 1 a} {sumi a 1} {sumi 1 1 b}")
	key := kb.BuildKey(mockContext())
	assert.Equal(t, "<ARGN> <BAD-TYPE> <BAD-TYPE> <BAD-TYPE>", key)
}

func TestMaxMin(t *testing.T) {
	testExpression(t, mockContext(), `{maxi 1 1} {maxi 1 2} {maxi 5 1} {mini 1 1} {mini 1 2} {mini 5 1}`, "1 2 5 1 1 1")
}

func TestArithmaticf(t *testing.T) {
	kb, _ := NewStdKeyBuilder().Compile("{sumf {1} {4}} {multf {1} 2} {divf {1} 2} {subf {1} 10}")
	key := kb.BuildKey(mockContext("ab", "100", "1000000", "5000000.123456", "22"))
	assert.Equal(t, "122 200 50 90", key)
}
func TestArithmaticfError(t *testing.T) {
	kb, _ := NewStdKeyBuilder().Compile("{sumf 1} {sumf 1 a} {sumf a 1} {sumf 1 2 a}")
	key := kb.BuildKey(mockContext())
	assert.Equal(t, "<ARGN> <BAD-TYPE> <BAD-TYPE> <BAD-TYPE>", key)
}