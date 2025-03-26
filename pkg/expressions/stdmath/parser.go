package stdmath

import (
	"errors"
	"strconv"
)

type (
	Expr interface {
		Eval(ctx Context) float64
		ToFunction() func(ctx Context) float64 // TODO: Is this more performant??? (less interfaces, more closures)
	}
	exprVal struct {
		v float64
	}
	exprNamedVar struct {
		name string
	}
	exprIntVar struct {
		idx int
	}
	exprUnary struct {
		op OpUnary
		ex Expr
	}
	exprBinary struct {
		op          OpFunc
		opCode      string
		left, right Expr
	}
)

func (s *exprVal) Eval(ctx Context) float64 {
	return s.v
}
func (s *exprNamedVar) Eval(ctx Context) float64 {
	return ctx.GetKey(s.name)
}
func (s *exprIntVar) Eval(ctx Context) float64 {
	return ctx.GetMatch(s.idx)
}
func (s *exprUnary) Eval(ctx Context) float64 {
	return s.op(s.ex.Eval(ctx))
}
func (s *exprBinary) Eval(ctx Context) float64 {
	return s.op(s.left.Eval(ctx), s.right.Eval(ctx))
}

func (s *exprVal) ToFunction() func(ctx Context) float64 {
	return func(ctx Context) float64 {
		return s.v
	}
}
func (s *exprNamedVar) ToFunction() func(ctx Context) float64 {
	return func(ctx Context) float64 {
		return ctx.GetKey(s.name)
	}
}
func (s *exprIntVar) ToFunction() func(ctx Context) float64 {
	return func(ctx Context) float64 {
		return ctx.GetMatch(s.idx)
	}
}
func (s *exprUnary) ToFunction() func(ctx Context) float64 {
	sub := s.ex.ToFunction()
	return func(ctx Context) float64 {
		return s.op(sub(ctx))
	}
}
func (s *exprBinary) ToFunction() func(ctx Context) float64 {
	left := s.left.ToFunction()
	right := s.right.ToFunction()
	return func(ctx Context) float64 {
		return s.op(left(ctx), right(ctx))
	}
}

func Compile(expr string) (Expr, error) {
	tokens, err := tokenizeExpr(expr)
	if err != nil {
		return nil, err
	}

	scanner := tokenScanner{tokens}

	// TODO: Reduce? (static analysis like expressions?) But unlike expressions, we can actually analyze what it is

	return scanner.compileTokens()
}

type tokenScanner struct {
	next []token
}

func (s *tokenScanner) compileTokens() (Expr, error) {
	var ret *exprBinary

	eFirst, _ := s.getNextExpr()
	if s.done() {
		return eFirst, nil
	}

	for !s.done() {
		// TODO: Err check
		op, opCode, _ := s.getNextOp()
		nextExpr, _ := s.getNextExpr()

		switch {
		case ret == nil:
			// TODO: Errors
			// TODO: Move this outside of loop?
			ret = &exprBinary{
				left:   eFirst,
				op:     op,
				opCode: opCode,
				right:  nextExpr,
			}
		case isOpBefore(ret.opCode, opCode): // eg. 3*3+3
			ret = &exprBinary{
				left:   ret,
				op:     op,
				opCode: opCode,
				right:  nextExpr,
			}
		default: // eg 3+3*3
			ret.right = &exprBinary{
				left:   ret.right,
				op:     op,
				opCode: opCode,
				right:  nextExpr,
			}
		}
	}

	return ret, nil
}

func (s *tokenScanner) getNextExpr() (Expr, error) {
	token := s.pop()
	switch token.t {
	case typeLiteral, typeGroup:
		return compileToken(token)

	case typeMod:
		modifier := uniOps[token.val]
		next, err := s.getNextExpr()
		if err != nil {
			return nil, err
		}

		return &exprUnary{
			op: modifier,
			ex: next,
		}, nil

	default:
		return nil, errors.New("unexpected token")
	}
}

func (s *tokenScanner) getNextOp() (OpFunc, string, error) {
	switch s.peek().t {
	case typeOp:
		token := s.pop()
		op := ops[token.val] // TODO: Erro check
		return op, token.val, nil
	case typeGroup: // special case, implied multiplication
		return ops["*"], "*", nil
	default:
		return nil, "", errors.New("expected operation")
	}
}

func (s *tokenScanner) pop() token {
	ret := s.next[0]
	s.next = s.next[1:]
	return ret
}

func (s *tokenScanner) peek() token {
	return s.next[0]
}

func (s *tokenScanner) done() bool {
	return len(s.next) == 0
}

// Turn a single token (literal or group) into an expression
func compileToken(t token) (Expr, error) {
	switch {
	case t.t == typeLiteral && isBraceBoxed(t.val):
		inner := t.val[1 : len(t.val)-1]
		if idx, err := strconv.Atoi(inner); err == nil {
			return &exprIntVar{idx}, nil
		}
		return &exprNamedVar{inner}, nil
	case t.t == typeLiteral:
		if v, err := strconv.ParseFloat(t.val, 64); err == nil {
			return &exprVal{v}, nil
		}
		return &exprNamedVar{t.val}, nil
	case t.t == typeGroup:
		return Compile(t.val)
	}

	return nil, errors.New("unexpected type")
}

// String like "{xxx}"
func isBraceBoxed(s string) bool {
	return len(s) > 0 && s[0] == '{' && s[len(s)-1] == '}'
}
