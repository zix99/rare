package stdmath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleEval(t *testing.T) {
	testFormula(t, mockContext(), "2*3", 6.0)
	testFormula(t, mockContext(), "2+3", 5.0)
	testFormula(t, mockContext(), "2-3.5", -1.5)
	testFormula(t, mockContext(), "10/2", 5.0)
	testFormula(t, mockContext(), "500*10", 5000.0)
}

func TestSimpleOrderOfOps(t *testing.T) {
	testFormula(t, mockContext("x", 123.0), "x*2+2", 248.0)
	testFormula(t, mockContext("x", 123.0), "2+x*2", 248.0)
	testFormula(t, mockContext("x", 123.0), "2+2*x", 248.0)
}

func TestParensFormula(t *testing.T) {
	ctx := mockContext("x", 5.0, "y", 12.0)
	testFormula(t, ctx, "x*(y+2)", 70.0)
	testFormula(t, ctx, "x*((y+2)/2)", 35.0)
	testFormula(t, ctx, "x*(y+2/2)", 5.0*13.0)
}

func TestNegativeNumbers(t *testing.T) {
	ctx := mockContext("x", 5.0)
	testFormula(t, ctx, "5 + -2", 3.0)
	testFormula(t, ctx, "8 + -x", 3.0)
	testFormula(t, ctx, "2 + -(3-2)", 1.0)
}

func TestMoreComplex(t *testing.T) {
	testFormula(t, nil, "cos(3.1415926535)", -1.0)
}

func TestImpliedMultiplication(t *testing.T) {
	testFormula(t, nil, "3(2)", 6.0)
	testFormula(t, nil, "1+3(2)", 7.0)
}

func TestComparisons(t *testing.T) {
	testFormula(t, nil, "1 <= 2", 1.0)
	testFormula(t, nil, "1 >= 2", 0.0)
}

func TestExplicitVariable(t *testing.T) {
	testFormula(t, mockContext("x", 150.0), "{x}/50", 3.0)
	testFormula(t, mockContext(), "{1}+3.0", 3.0)
}

func mockContext(eles ...interface{}) Context {
	m := make(map[string]float64)
	for i := 0; i < len(eles); i += 2 {
		m[eles[i].(string)] = eles[i+1].(float64)
	}
	return &SimpleContext{namedVals: m}
}

func testFormula(t *testing.T, ctx Context, f string, expected float64) {
	t.Run(f, func(t *testing.T) {
		expr, err := Compile(f)
		assert.NoError(t, err)

		ret := expr.Eval(ctx)
		assert.Equal(t, expected, ret)
	})
}

// BenchmarkFormula-4   	25900489	        42.30 ns/op	       0 B/op	       0 allocs/op
func BenchmarkFormula(b *testing.B) {
	expr, _ := Compile("2 + 5 * x")
	ctx := mockContext("x", 5.0)
	f := expr.ToFunction()
	for range b.N {
		// expr.Eval(ctx)
		f(ctx)
	}
}
