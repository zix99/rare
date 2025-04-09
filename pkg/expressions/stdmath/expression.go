package stdmath

type (
	Expr interface {
		Eval(ctx Context) float64
	}
	exprVal struct {
		v float64
	}
	exprNamedVar struct {
		name string
	}
	exprIndexVar struct {
		idx int
	}
	exprUnary struct {
		op OpUnary
		ex Expr
	}
	exprBinary struct {
		op          OpFunc
		opCode      OpCode
		left, right Expr
	}
)

func (s *exprVal) Eval(ctx Context) float64 {
	return s.v
}
func (s *exprNamedVar) Eval(ctx Context) float64 {
	return ctx.GetKey(s.name)
}
func (s *exprIndexVar) Eval(ctx Context) float64 {
	return ctx.GetMatch(s.idx)
}
func (s *exprUnary) Eval(ctx Context) float64 {
	return s.op(s.ex.Eval(ctx))
}
func (s *exprBinary) Eval(ctx Context) float64 {
	return s.op(s.left.Eval(ctx), s.right.Eval(ctx))
}
