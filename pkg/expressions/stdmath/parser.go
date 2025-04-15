package stdmath

import (
	"regexp"
	"strconv"
)

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
		return nil, ErrUnexpectedEnd
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

		switch opCodeOrder(lastOpCode, peekOp) {
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
				left:   simplify(ret),
				op:     op,
				opCode: opCode,
				right:  simplify(expr),
			}
		}
	}

	return simplify(ret), nil
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
		return nil, ErrExpectedExpression
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
			return nil, "", ErrUnknownOperation
		}
		return op, OpCode(token.val), nil
	case typeGroup: // special case, implied multiplication
		return ops["*"], "*", nil
	default:
		return nil, "", ErrExpectedOperation
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
	case t.t == typeLiteral && isBoxed(t.val):
		inner := t.val[1 : len(t.val)-1]
		if idx, err := strconv.Atoi(inner); err == nil {
			return &exprIndexVar{idx}, nil
		}
		return &exprNamedVar{inner}, nil

	case t.t == typeLiteral: // numeric literal
		// 0b, 0x, or 10 int
		if v, err := strconv.ParseInt(t.val, 0, 64); err == nil {
			return &exprVal{float64(v)}, nil
		}

		// const float
		if v, err := strconv.ParseFloat(t.val, 64); err == nil {
			return &exprVal{v}, nil
		}

		if !validVariableName(t.val) {
			return nil, ErrTokenizerNumeric
		}

		return &exprNamedVar{t.val}, nil

	case t.t == typeGroup:
		return Compile(t.val)
	}

	return nil, ErrExpectedExpression
}

// String like "{xxx}"
func isBoxed(s string) bool {
	return len(s) >= 2 && s[0] == '[' && s[len(s)-1] == ']'
}

var validVariableRegex = regexp.MustCompile("(?i)^[a-z][a-z0-9]*$")

// Check for valid variable names (in implicit cases)
func validVariableName(s string) bool {
	return validVariableRegex.MatchString(s)
}
