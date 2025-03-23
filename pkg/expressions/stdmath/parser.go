package stdmath

import (
	"iter"
	"math"
	"slices"
	"strconv"
	"strings"
)

type Operation rune

type Expr struct {
	Op          OpFunc
	opCode      string
	left, right *Expr
	value       *float64
	named       *string
}

type OpFunc func(left, right float64) float64

var orderOfOps = []string{"*", "/", "+", "-"}

func isOpBefore(op0, op1 string) bool {
	for _, op := range orderOfOps {
		if op == op0 { // saw op0 first
			return true
		}
		if op == op1 { // saw op1 first
			return false
		}
	}
	panic("fixme")
}

var ops = map[string]OpFunc{
	"+": func(left, right float64) float64 { return left + right },
	"*": func(left, right float64) float64 { return left * right },
	"-": func(left, right float64) float64 { return left - right },
	"/": func(left, right float64) float64 { return left / right },
	"^": func(left, right float64) float64 { return math.Pow(left, right) },
}

// func Compile(s string) expressions.KeyBuilderStage {
// 	return func(kbc expressions.KeyBuilderContext) string {

// 	}
// }

type Context interface {
	GetMatch(int) float64
	GetKey(string) float64
}

type SimpleContext struct {
	namedVals map[string]float64
}

func (s *SimpleContext) GetMatch(idx int) float64 {
	return 0
}

func (s *SimpleContext) GetKey(k string) float64 {
	return s.namedVals[k]
}

func Compile(expr string) *Expr {
	return CompileEx(slices.Collect(tokenizeExpr(expr))...)
}

func CompileEx(tokens ...string) *Expr {
	if len(tokens) == 1 {
		tok := tokens[0]
		if val, err := strconv.ParseFloat(tok, 64); err == nil {
			return &Expr{
				value: &val,
			}
		} else {
			// Assume expression
			return &Expr{
				named: &tok,
			}
		}
	}

	var ret *Expr

	// lastOp := tokens[i+1]
	for i := 0; i <= len(tokens)-3; i += 2 {
		top := tokens[i+1]
		tright := tokens[i+2]

		// if op is higher priority, put it in a group and nest
		if ret == nil {
			tleft := tokens[i]
			ret = &Expr{
				left:   CompileEx(tleft),
				Op:     ops[top],
				opCode: top,
				right:  CompileEx(tright),
			}
		} else if isOpBefore(ret.opCode, top) { // eg. 3*3+3
			ret = &Expr{
				left:   ret,
				Op:     ops[top],
				opCode: top,
				right:  CompileEx(tright),
			}
		} else { // eg 3+<3*3>
			ret.right = &Expr{
				left:   ret.right,
				Op:     ops[top],
				opCode: top,
				right:  CompileEx(tright),
			}
		}
	}

	return ret
}

func (s *Expr) Eval(ctx Context) float64 {
	if s.value != nil {
		return *s.value
	}
	if s.named != nil {
		return ctx.GetKey(*s.named)
	}

	return s.Op(s.left.Eval(ctx), s.right.Eval(ctx))
}

type tokenType int

const (
	literal tokenType = 1 << iota
	grouping
	operation
)

type token struct {
	val string
	op  tokenType
}

func tokenizeExpr(s string) iter.Seq[string] {
	return func(yield func(string) bool) {
		var token strings.Builder
		parens := 0

		for _, r := range s {
			switch {
			case r == '(':
				if parens > 0 {
					token.WriteRune('(')
				}
				parens++
			case r == ')':
				parens--
				if parens == 0 {
					if !yield(token.String()) {
						return
					}
					token.Reset()
				} else if parens < 0 {
					// error
					panic("fixme")
				} else {
					token.WriteRune(')')
				}
			case r == ' ':
				// skip
			case parens == 0 && in(r, '+', '-', '*', '/', '^'): // operation FIXME: Use actual ops
				if token.Len() > 0 && !yield(token.String()) {
					return
				}
				if !yield(string(r)) {
					return
				}
				token.Reset()
			default: // token/value
				token.WriteRune(r)
			}
		}

		if token.Len() > 0 {
			yield(token.String())
		}
	}
}

func in[T comparable](s T, eles ...T) bool {
	return slices.Contains(eles, s)
}
