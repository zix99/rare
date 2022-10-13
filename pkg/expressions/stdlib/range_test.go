package stdlib

import "testing"

func TestArrayMap(t *testing.T) {
	testExpression(
		t,
		mockContext(mockArray("joe", "is", "cool")),
		`{$join {$map {0} "{0}bob"} ", "}`,
		"joebob, isbob, coolbob",
	)
}

func TestArrayReduce(t *testing.T) {
	testExpression(
		t,
		mockContext("0 1 2 5"),
		`{$reduce {$split {0} " "} "{sumi {0} {1}}"}`,
		"8",
	)
}

func TestArraySlice(t *testing.T) {
	testExpression(
		t,
		mockContext("0 1 2 5"),
		`{$join {$slice {$split {0} " "} 1 2}}`,
		"1 2",
	)
}

func TestArrayFilter(t *testing.T) {
	testExpression(
		t,
		mockContext(mockArray("a", "123", "b", "455")),
		`{$join {$filter {0} "{isnum {0}}"}}`,
		"123 455",
	)
}
