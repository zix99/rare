package expressions

import "testing"

func TestBaseName(t *testing.T) {
	testExpression(t, mockContext("ab/c/d"), "{basename {0}} {basename a b}", "d <ARGN>")
}

func TestDirName(t *testing.T) {
	testExpression(t, mockContext("ab/c/d"), "{dirname {0}}", "ab/c")
}

func TestExtName(t *testing.T) {
	testExpression(t, mockContext(), "{extname a/b/c} {extname a/b/c.jpg}", " .jpg")
}
