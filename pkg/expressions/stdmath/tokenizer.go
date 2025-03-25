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
		case r == '(':
			if parens > 0 {
				sb.WriteRune('(')
			} else if sb.Len() > 0 {
				// has previous token
				ret = append(ret, token{sb.String(), typeMod})
				sb.Reset()
			}
			parens++
		case r == ')':
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
		case r == ' ':
			// skip
		case parens == 0 && sb.Len() == 0 && (len(ret) == 0 || (len(ret) > 0 && ret[len(ret)-1].t == typeOp)) && in(r, '-'): // modifier (unary)
			ret = append(ret, token{string(r), typeMod})
			sb.Reset()
		case parens == 0 && in(r, '+', '-', '*', '/', '^'): // operation FIXME: Use actual ops
			if sb.Len() > 0 {
				ret = append(ret, token{sb.String(), typeLiteral})
			}
			ret = append(ret, token{string(r), typeOp})
			sb.Reset()
		default: // token/value
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
