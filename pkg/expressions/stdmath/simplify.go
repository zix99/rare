package stdmath

/*
This is a better-than-nothing simplifier
that looks at an expression, checks if the context is accessed
and assumes it's a constant if not.

There are significantly better ways to do this, and many cases
that this won't catch (eg. 2+x+2 will look like both nodes
use `x`)
*/

type simplifyContext struct {
	hits int
}

func (s *simplifyContext) GetMatch(idx int) float64 {
	s.hits++
	return 0
}

func (s *simplifyContext) GetKey(name string) float64 {
	s.hits++
	return 0
}

func simplify(expr Expr) Expr {
	ctx := &simplifyContext{}
	if val := expr.Eval(ctx); ctx.hits == 0 {
		return &exprVal{val}
	}
	return expr
}
