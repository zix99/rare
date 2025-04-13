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

func TestOpsBasic(t *testing.T) {
	testOp(t, "1+2.5", 3.5)
	testOp(t, "1-3", -2)
	testOp(t, "2*5", 10)
	testOp(t, "12/6", 2)
	testOp(t, "2^3", 8)
	testOp(t, "5 % 2", 1)
}

func TestOpsShift(t *testing.T) {
	testOp(t, "1<<2", 4)
	testOp(t, "4>>1", 2)
	testOp(t, "5 & 0b100", 4)
	testOp(t, "0b10 | 0b01", 3)
}

func TestOpCompare(t *testing.T) {
	testOp(t, "1 && 0", 0)
	testOp(t, "1 || 0", 1)
}

func TestUnaryBasic(t *testing.T) {
	testOp(t, "!1", 0.0)
	testOp(t, "!(1 > 2)", 1.0)
	testOp(t, "abs(-3)", 3)
	testOp(t, "-1", -1)
}

func TestTrig(t *testing.T) {
	testOp(t, "sin(0)", 0)
	testOp(t, "asin(0)", 0)
	testOp(t, "cos(0)", 1)
	testOp(t, "acos(0)", 1.5707963267948966)
	testOp(t, "tan(0)", 0)
	testOp(t, "atan(0)", 0)
}

func TestRounding(t *testing.T) {
	testOp(t, "round(3.5)", 4)
	testOp(t, "floor(3.5)", 3)
	testOp(t, "ceil(3.5)", 4)
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
