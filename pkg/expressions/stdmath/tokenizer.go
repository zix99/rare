package stdmath

import (
	"errors"
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
	ErrTokenizerOverclosed = errors.New("over-closed parenthesis")
	ErrTokenizerUnclosed   = errors.New("unclosed parenthesis")
)

func tokenizeExpr(s string) ([]token, error) {
	ret := make([]token, 0)

	var sb strings.Builder
	parens := 0

	for i := 0; i < len(s); i++ {
		r := s[i]

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

			if _, uniOk := uniOps[OpCode(prev)]; uniOk {
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
		case parens == 0 && prefixInOps(s[i:]) != nil:
			if sb.Len() > 0 {
				ret = append(ret, token{sb.String(), typeLiteral})
				sb.Reset()
			}

			opCode := *prefixInOps(s[i:])
			ret = append(ret, token{string(opCode), typeOp})
			i += len(opCode) - 1

		// Token continuation
		default:
			sb.WriteByte(r)
		}
	}

	if parens > 0 {
		return nil, ErrTokenizerUnclosed
	}

	if sb.Len() > 0 {
		ret = append(ret, token{sb.String(), typeLiteral})
	}

	return ret, nil
}
