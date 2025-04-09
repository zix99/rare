package stdmath

import (
	"maps"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderOfOperationExhaustive(t *testing.T) {
	allOps := slices.Collect(maps.Keys(ops))
	flattenedOrder := slices.Concat(orderOfOps...)
	assert.ElementsMatch(t, allOps, flattenedOrder)
}

func TestOpCodes(t *testing.T) {
	assert.Equal(t, opCodeOrder("*", "+"), -1)
	assert.Equal(t, opCodeOrder("-", "+"), 0)
	assert.Equal(t, opCodeOrder("*", "^"), 1)
}

func TestPrefixInOps(t *testing.T) {
	assert.EqualValues(t, "*", *prefixInOps("*"))
	assert.EqualValues(t, "*", *prefixInOps("*b"))
	assert.EqualValues(t, "<=", *prefixInOps("<= b"))
	assert.EqualValues(t, "<", *prefixInOps("< b"))
}

func TestUnaryNot(t *testing.T) {
	testFormula(t, nil, "!1", 0.0)
	testOp(t, "!(1 > 2)", 1.0)
}

func testOp(t *testing.T, f string, expected float64) {
	t.Helper()

	expr, err := Compile(f)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	ret := expr.Eval(nil)
	if !assert.Equal(t, expected, ret) {
		debugWriteTree(expr, 0)
	}
}
