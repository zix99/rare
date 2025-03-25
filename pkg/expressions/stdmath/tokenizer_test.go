package stdmath

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenizeExpression(t *testing.T) {
	testTokenizer(t, "123+1", "123(0) +(2) 1(0)")
	testTokenizer(t, "2-3.5", "2(0) -(2) 3.5(0)")
}

func TestTokenizeParens(t *testing.T) {
	testTokenizer(t, "1 + (1+1)", "1(0) +(2) 1+1(1)")
	testTokenizer(t, "1+( 1 +1)", "1(0) +(2) 1+1(1)")
	testTokenizer(t, "10+20", "10(0) +(2) 20(0)")
	testTokenizer(t, "1 + (1+(2*3))*5", "1(0) +(2) 1+(2*3)(1) *(2) 5(0)")
}

func TestTokenizeUnary(t *testing.T) {
	testTokenizer(t, "-x", "-(3) x(0)")
	testTokenizer(t, "abs(x)", "abs(3) x(1)")
	testTokenizer(t, "1+ -x", "1(0) +(2) -(3) x(0)")
	testTokenizer(t, "1 + abs(2)", "1(0) +(2) abs(3) 2(1)")
	testTokenizer(t, "1 + abs(-2)", "1(0) +(2) abs(3) -2(1)")
	testTokenizer(t, "1 + abs(x+3)", "1(0) +(2) abs(3) x+3(1)")
	testTokenizer(t, "2 + -(3-2)", "2(0) +(2) -(3) 3-2(1)")
}

func TestTokenizerVariables(t *testing.T) {
	testTokenizer(t, "abc", "abc(0)")
}

// test forumla parses into expects
// expects is a stringified token result in the format "val(type) ..."
func testTokenizer(t *testing.T, formula, expects string) {
	t.Run(formula, func(t *testing.T) {
		tokens, err := tokenizeExpr(formula)
		assert.NoError(t, err)
		assert.NotNil(t, tokens)

		// stringify to make it easier to test
		var sb strings.Builder
		for _, ele := range tokens {
			if sb.Len() > 0 {
				sb.WriteRune(' ')
			}
			sb.WriteString(fmt.Sprintf("%s(%d)", ele.val, ele.t))
		}

		// assert
		assert.Equal(t, expects, sb.String())
	})
}
