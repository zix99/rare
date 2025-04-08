package stdmath

import (
	"errors"
	"strconv"
)

/* TODO
- Reduce static expressions for pre-compute (Simplify)
- Support generic types? Anything numeric should be able to compute generically
*/

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
	exprIntVar struct {
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
func (s *exprIntVar) Eval(ctx Context) float64 {
	return ctx.GetMatch(s.idx)
}
func (s *exprUnary) Eval(ctx Context) float64 {
	return s.op(s.ex.Eval(ctx))
}
func (s *exprBinary) Eval(ctx Context) float64 {
	return s.op(s.left.Eval(ctx), s.right.Eval(ctx))
}

func Compile(expr string) (Expr, error) {
	tokens, err := tokenizeExpr(expr)
	if err != nil {
		return nil, err
	}

	scanner := tokenScanner{expr, tokens}

	return scanner.compileTokens("")
}

type tokenScanner struct {
	s    string
	next []token
}

func (s *tokenScanner) compileTokens(lastOpCode OpCode) (ret Expr, err error) {
	if s.done() {
		return nil, errors.New("unexpected end")
	}

	ret, err = s.getNextExpr()
	if err != nil {
		return nil, err
	}

	for !s.done() {
		_, peekOp, err := s.getNextOp(false)
		if err != nil {
			return nil, err
		}
		order := opCodeOrder(lastOpCode, peekOp)

		switch order {
		case -1: // * -> +
			return ret, nil // no op, just expression
		case 0: // + -> +-
			return ret, nil
		case 1: // + -> *
			// recurse
			op, opCode, err := s.getNextOp(true)
			if err != nil {
				return nil, err
			}

			expr, err := s.compileTokens(opCode)
			if err != nil {
				return nil, err
			}
			ret = &exprBinary{
				left:   ret,
				op:     op,
				opCode: opCode,
				right:  expr,
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
		modifier := uniOps[OpCode(token.val)]
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

func (s *tokenScanner) getNextOp(pop bool) (OpFunc, OpCode, error) {
	switch s.peek().t {
	case typeOp:
		token := s.peek()
		if pop {
			s.pop()
		}
		op, ok := ops[OpCode(token.val)]
		if !ok {
			return nil, "", errors.New("unrecognized op")
		}
		return op, OpCode(token.val), nil
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
