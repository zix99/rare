package stdmath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimplify(t *testing.T) {
	s := simplify(&exprUnary{uniOps["-"], &exprVal{2.0}})
	assert.IsType(t, &exprVal{}, s)
	assert.Equal(t, s.Eval(nil), -2.0)

	s = simplify(&exprBinary{
		left:  &exprVal{1.0},
		op:    ops["+"],
		right: &exprVal{2.0},
	})
	assert.IsType(t, &exprVal{}, s)
	assert.Equal(t, s.Eval(nil), 3.0)

	s = simplify(&exprUnary{uniOps["-"], &exprIndexVar{1}})
	assert.IsType(t, &exprUnary{}, s)
}
