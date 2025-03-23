package stdmath

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenizeExpression(t *testing.T) {
	tok := slices.Collect(tokenizeExpr("123 +1"))
	assert.Equal(t, []string{"123", "+", "1"}, tok)
}

func TestTokenizeParens(t *testing.T) {
	assert.Equal(t, []string{"1", "+", "1+1"}, slices.Collect(tokenizeExpr("1 + (1+1)")))
	assert.Equal(t, []string{"10", "+", "20"}, slices.Collect(tokenizeExpr("10+20")))
	assert.Equal(t, []string{"1", "+", "1+(2*3)", "*", "5"}, slices.Collect(tokenizeExpr("1 + (1+(2*3))*5")))
}

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

func mockContext(eles ...interface{}) Context {
	m := make(map[string]float64)
	for i := 0; i < len(eles); i += 2 {
		m[eles[i].(string)] = eles[i+1].(float64)
	}
	return &SimpleContext{namedVals: m}
}

func testFormula(t *testing.T, ctx Context, f string, expected float64) {
	t.Run(f, func(t *testing.T) {
		expr := Compile(f)
		ret := expr.Eval(ctx)
		assert.Equal(t, expected, ret)
	})
}
