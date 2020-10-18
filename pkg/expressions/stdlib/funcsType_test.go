package stdlib

import "testing"

func TestIsInt(t *testing.T) {
	testExpression(t, mockContext(), "{isint 123} {isint 123.0} {isint abc}", "1  ")
}

func TestIsNum(t *testing.T) {
	testExpression(t, mockContext(), "{isnum 123} {isnum 123.0} {isnum abc}", "1 1 ")
}
