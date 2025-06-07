package stdlib

import (
	"testing"

	"github.com/zix99/rare/pkg/expressions"
)

// Focus on testing the math integration
// tests cover math itself
func TestMath(t *testing.T) {
	testExpression(t, mockContext("25"), "{! 2*[0]}", "50")
	testExpression(t, mockContext("25"), "{! 2 * [0]}", "50")
	testExpression(t, mockContext("25"), "{! [0] / 5}", "5")
	testExpression(t, mockContext("abc"), "{! [0]*2}", "<BAD-TYPE>")

	testExpression(t, &expressions.KeyBuilderContextArray{
		Keys: map[string]string{
			"val": "5",
			"x":   "2",
		},
	}, "{! val*[x]}", "10")

	testExpressionErr(t, mockContext("1"), "{! 1+{0}}", "<CONST>", ErrConst) // Need to use bracket notation
	testExpressionErr(t, mockContext(), "{! 1+}", "<PARSE-ERROR>", ErrParsing)
}
