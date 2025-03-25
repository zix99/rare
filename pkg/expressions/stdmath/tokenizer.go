package stdmath

import (
	"errors"
	"slices"
	"strings"
)

type tokenType int

const (
	typeLiteral tokenType = iota
	typeGroup
	typeOp
	typeMod // unary
)

type token struct {
	val string
	t   tokenType
}

var (
	ErrTokenizerOverclosed = errors.New("over-closed paranthesis")
)

func tokenizeExpr(s string) ([]token, error) {
	ret := make([]token, 0)

	var sb strings.Builder
	parens := 0

	for _, r := range s {
		switch {
		// parens management
		case r == '(' && parens > 0:
			// Nested paren
			sb.WriteRune('(')
			parens++
		case r == '(' && sb.Len() > 0:
			// previous token. Possibly unary op or implicit multiply (literal)
			prev := sb.String()
			sb.Reset()

			if _, uniOk := uniOps[prev]; uniOk {
				ret = append(ret, token{prev, typeMod})
			} else {
				ret = append(ret, token{prev, typeLiteral})
			}
			parens++
		case r == '(':
			// Other paren
			parens++
		case r == ')': // end paren
			parens--
			if parens == 0 {
				ret = append(ret, token{sb.String(), typeGroup})
				sb.Reset()
			} else if parens < 0 {
				// error
				return nil, ErrTokenizerOverclosed
			} else {
				sb.WriteRune(')')
			}

		// Skip whitespace
		case r == ' ':
			// skip

		// negative unary op on literal or group
		case parens == 0 && sb.Len() == 0 && (len(ret) == 0 || (len(ret) > 0 && ret[len(ret)-1].t == typeOp)) && r == '-':
			ret = append(ret, token{string(r), typeMod})

		// operator
		case parens == 0 && in(r, '+', '-', '*', '/', '^'): // operation FIXME: Use actual ops
			if sb.Len() > 0 {
				ret = append(ret, token{sb.String(), typeLiteral})
			}
			ret = append(ret, token{string(r), typeOp})
			sb.Reset()

		// Token continuation
		default:
			sb.WriteRune(r)
		}
	}

	if sb.Len() > 0 {
		ret = append(ret, token{sb.String(), typeLiteral})
	}

	return ret, nil
}

func in[T comparable](s T, eles ...T) bool {
	return slices.Contains(eles, s)
}
