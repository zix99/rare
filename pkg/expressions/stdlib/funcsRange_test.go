package stdlib

import (
	"rare/pkg/expressions"
	"testing"
)

func TestArray(t *testing.T) {
	testExpression(t, mockContext("q"), "{$ {0} {1} 22}", "q\x00\x0022")
	testExpression(t, mockContext("q"), `{$ "{0} hi" 22}`, "q hi\x0022")
	testExpression(t, mockContext("q"), "{$ {0}}", "q")
}

func TestArrayLen(t *testing.T) {
	testExpression(t, mockContext("abc"), "{@len {0}}", "1")
	testExpression(t, mockContext(expressions.MakeArray("a", "bc", "c")), "{@len {0}}", "3")
	testExpression(t, mockContext(""), "{@len {0}}", "0")
	testExpressionErr(t, mockContext(), "{@len a b}", "<ARGN>", ErrArgCount)
}

func TestArraySplit(t *testing.T) {
	testExpression(
		t,
		mockContext("a\tb\tc"),
		`{@split {0} "\t"}`,
		"a\x00b\x00c",
	)
	testExpression(
		t,
		mockContext("abc"),
		`{@split {0} "\t"}`,
		"abc",
	)
	testExpression(
		t,
		mockContext("a b\tc"),
		`{@split {0}}`,
		"a\x00b\tc",
	)
	testExpressionErr(
		t,
		mockContext("a b\tc"),
		`{@split {0} ""}`,
		"<EMPTY>",
		ErrEmpty,
	)
	testExpressionErr(
		t,
		mockContext("a b\tc"),
		`{@split {0} "" "c"}`,
		"<ARGN>",
		ErrArgCount,
	)
}

func TestArrayJoin(t *testing.T) {
	testExpression(
		t,
		mockContext("a\x00b\x00c"),
		`{@join {0} ", "}`,
		"a, b, c",
	)
	testExpression(
		t,
		mockContext("a\x00b\x00c"),
		`{@join {0} ""}`,
		"abc",
	)
	testExpression(
		t,
		mockContext("a"),
		`{@join {0} ", "}`,
		"a",
	)
	testExpression(
		t,
		mockContext("a\x00b\x00c"),
		`{@join {0}}`,
		"a b c",
	)
	testExpressionErr(
		t,
		mockContext("a\x00b\x00c"),
		`{@join {0} ", " "c"}`,
		"<ARGN>",
		ErrArgCount,
	)
}

func TestArraySelect(t *testing.T) {
	testExpression(t, mockContext(expressions.MakeArray("aq", "bc", "c")), `{@select {0} 0}`, "aq")
	testExpression(t, mockContext(expressions.MakeArray("aq", "bc", "c")), `{@select {0} 1}`, "bc")
	testExpression(t, mockContext(expressions.MakeArray("aq", "bc", "c")), `{@select {0} 2}`, "c")
	testExpression(t, mockContext(expressions.MakeArray("aq", "bc", "c")), `{@select {0} 3}`, "")
	testExpression(t, mockContext(expressions.MakeArray("aq", "bc", "c")), `{@select {0} -1}`, "c")
}

func TestArrayMap(t *testing.T) {
	testExpression(
		t,
		mockContext(expressions.MakeArray("joe", "is", "cool")),
		`{@join {@map {0} "{0}bob"} ", "}`,
		"joebob, isbob, coolbob",
	)
	testExpression(
		t,
		mockContext(expressions.MakeArray("5", "1", "3")),
		`{@join {@map {0} "{multi {0} 2}"} ", "}`,
		"10, 2, 6",
	)
	testExpressionErr(
		t,
		mockContext(expressions.MakeArray("5", "1", "3")),
		`{@join {@map {0} "{multi {0} 2}" ""} ", "}`,
		"<ARGN>",
		ErrArgCount,
	)
}

func TestArrayReduce(t *testing.T) {
	testExpression(
		t,
		mockContext("0 1 2 5"),
		`{@reduce {@split {0} " "} "{sumi {0} {1}}"}`,
		"8",
	)
	testExpressionErr(
		t,
		mockContext("0 1 2 5"),
		`{@reduce {@split {0} " "} "{sumi {0} {1}}" bla 2}`,
		"<ARGN>",
		ErrArgCount,
	)
	testExpression(
		t,
		mockContext(expressions.MakeArray("1", "1", "3", "5")),
		`{@reduce {0} "{sumi {0} {1}}"}`,
		"10",
	)

	// With initial
	testExpression(t,
		mockContext(expressions.MakeArray("2", "1", "3", "5")),
		`{@reduce {0} "{subi {0} {1}}" 0}`, "-11")
}

func TestArraySlice(t *testing.T) {
	testExpression(
		t,
		mockContext("0 1 2 5"),
		`{@join {@slice {@split {0} " "} 1 2}}`,
		"1 2",
	)
	testExpression(
		t,
		mockContext("0 1 2 5"),
		`{@join {@slice {@split {0} " "} 1}}`,
		"1 2 5",
	)
	testExpression(
		t,
		mockContext("0 1 2 5"),
		`{@join {@slice {@split {0} " "} 0 50}}`,
		"0 1 2 5",
	)
	testExpression(
		t,
		mockContext("0 1 2 5"),
		`{@join {@slice {@split {0} " "} 10 2}}`,
		"",
	)
	testExpression(
		t,
		mockContext("0 1 2 5"),
		`{@join {@slice {@split {0} " "} -3 2}}`,
		"1 2",
	)
	testExpressionErr(
		t,
		mockContext("0 1 2 5"),
		`{@join {@slice {@split {0} " "} 1 2 bla}}`,
		"<ARGN>",
		ErrArgCount,
	)
	testExpressionErr(
		t,
		mockContext("0 1 2 5"),
		`{@join {@slice {@split {0} " "} 1 bla}}`,
		"<CONST>",
		ErrConst,
	)
}

func TestArrayFilter(t *testing.T) {
	testExpression(
		t,
		mockContext(expressions.MakeArray("a", "123", "b", "455")),
		`{@join {@filter {0} "{isnum {0}}"}}`,
		"123 455",
	)
	testExpression(
		t,
		mockContext(expressions.MakeArray("a", "123", "b", "455")),
		`{@join {@filter {0} "1"}}`,
		"a 123 b 455",
	)
	testExpression(
		t,
		mockContext(expressions.MakeArray("a", "123", "b", "455")),
		`{@join {@filter {0} ""}}`,
		"",
	)
	testExpression( // filter with empty
		t,
		mockContext(expressions.MakeArray("", "123", "", "456")),
		`{@join {@filter {0} "1"} ","}`,
		",123,,456",
	)
	testExpressionErr(
		t,
		mockContext(expressions.MakeArray("a", "123", "b", "455")),
		`{@join {@filter {0}}}`,
		"<ARGN>",
		ErrArgCount,
	)
}

func TestArrayRange(t *testing.T) {
	// 1 Arg
	testExpression(t, mockContext("5"), "{@range {0}}", expressions.MakeArray("0", "1", "2", "3", "4"))
	testExpression(t, mockContext("0"), "{@range {0}}", expressions.MakeArray())
	testExpression(t, mockContext("-1"), "{@range {0}}", expressions.MakeArray("<VALUE>"))
	testExpression(t, mockContext("abc"), "{@range {0}}", expressions.MakeArray("<BAD-TYPE>"))

	// 2 Arg
	testExpression(t, mockContext("5"), "{@range 1 {0}}", expressions.MakeArray("1", "2", "3", "4"))
	testExpression(t, mockContext("0"), "{@range 0 {0}}", expressions.MakeArray())
	testExpression(t, mockContext("-1"), "{@range 0 {0}}", expressions.MakeArray("<VALUE>"))
	testExpression(t, mockContext("-1"), "{@range -1 2}", expressions.MakeArray("-1", "0", "1"))
	testExpression(t, mockContext("-1"), "{@range 5 3}", expressions.MakeArray("<VALUE>"))
	testExpression(t, mockContext("abc"), "{@range 0 {0}}", expressions.MakeArray("<BAD-TYPE>"))

	// 3 Arg
	testExpression(t, mockContext("5"), "{@range 1 {0} 1}", expressions.MakeArray("1", "2", "3", "4"))
	testExpression(t, mockContext("0"), "{@range 0 {0} 1}", expressions.MakeArray())
	testExpression(t, mockContext("-1"), "{@range 0 {0} 1}", expressions.MakeArray("<VALUE>"))
	testExpression(t, mockContext("-1"), "{@range -1 2 1}", expressions.MakeArray("-1", "0", "1"))
	testExpression(t, mockContext("-1"), "{@range 5 3 1}", expressions.MakeArray("<VALUE>"))
	testExpression(t, mockContext("abc"), "{@range 0 {0} 1}", expressions.MakeArray("<BAD-TYPE>"))
	testExpression(t, mockContext(), "{@range 5 1 -1}", expressions.MakeArray("5", "4", "3", "2"))

	// 4+ arg
	testExpressionErr(t, mockContext(), "{@range 1 2 3 4}", "<ARGN>", ErrArgCount)

	// Other error states
	testExpression(t, mockContext(), "{@range a}", "<BAD-TYPE>")
	testExpression(t, mockContext(), "{@range b 5}", "<BAD-TYPE>")
	testExpression(t, mockContext(), "{@range 0 5 c}", "<BAD-TYPE>")
	testExpression(t, mockContext(), "{@range 0 5 -1}", "<VALUE>")
	testExpression(t, mockContext(), "{@range 5 0 2}", "<VALUE>")
	testExpression(t, mockContext(), "{@range 0 5 0}", "<VALUE>")
}

func TestArrayFor(t *testing.T) {
	testExpressionErr(t, mockContext(), "{@for 5}", "<ARGN>", ErrArgCount)
	testExpression(t, mockContext(), "{@for 2 {lt {1} 5} {sumi {0} 2}}", expressions.MakeArray("2", "4", "6", "8", "10"))
}

func TestArrayIn(t *testing.T) {
	testExpression(t, mockContext("ab"), "{@in {0} {$ cd ab qef}}", "1")
	testExpression(t, mockContext("a"), "{@in {0} {$ cd ab qef}}", "")
	testExpression(t, mockContext("a"), `{@in {0} ""}`, "")

	testExpressionErr(t, mockContext("ab"), "{@in {0} {$ cd ab qef {1}}}", "<CONST>", ErrConst)
	testExpressionErr(t, mockContext("ab"), "{@in {0}}", "<ARGN>", ErrArgCount)
}

// BenchmarkRangeSum-4   	 4414395	       271.9 ns/op	       0 B/op	       0 allocs/op
func BenchmarkRangeSum(b *testing.B) {
	exp := NewStdKeyBuilder()
	ctx := mockContext(expressions.MakeArray("1", "1", "3", "5"))

	c, _ := exp.Compile("{@reduce {0} {sumi {0} {1}}}")

	for i := 0; i < b.N; i++ {
		c.BuildKey(ctx)
	}
}
