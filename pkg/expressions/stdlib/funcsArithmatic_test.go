package stdlib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArithmatic(t *testing.T) {
	kb, _ := NewStdKeyBuilder().Compile("{sumi {1} {4}} {multi {1} 2} {divi {1} 2} {subi {1} 10} {modi {1} 7}")
	key := kb.BuildKey(mockContext("ab", "100", "1000000", "5000000.123456", "22"))
	assert.Equal(t, "122 200 50 90 2", key)
}

func TestArithmaticError(t *testing.T) {
	kb, _ := NewStdKeyBuilder().Compile("{sumi 1} {sumi 1 a} {sumi a 1} {sumi 1 1 b}")
	key := kb.BuildKey(mockContext())
	assert.Equal(t, "<ARGN> <BAD-TYPE> <BAD-TYPE> <BAD-TYPE>", key)
}

func TestMaxMin(t *testing.T) {
	testExpression(t, mockContext(), `{maxi 1 1} {maxi 1 2} {maxi 5 1} {mini 1 1} {mini 1 2} {mini 5 1}`, "1 2 5 1 1 1")
}

func TestArithmaticf(t *testing.T) {
	kb, _ := NewStdKeyBuilder().Compile("{sumf {1} {4}} {multf {1} 2} {divf {1} 2} {subf {1} 10}")
	key := kb.BuildKey(mockContext("ab", "100", "1000000", "5000000.123456", "22"))
	assert.Equal(t, "122 200 50 90", key)
}
func TestArithmaticfError(t *testing.T) {
	kb, _ := NewStdKeyBuilder().Compile("{sumf 1} {sumf 1 a} {sumf a 1} {sumf 1 2 a}")
	key := kb.BuildKey(mockContext())
	assert.Equal(t, "<ARGN> <BAD-TYPE> <BAD-TYPE> <BAD-TYPE>", key)
}

func TestFloorCeilRound(t *testing.T) {
	testExpression(t, mockContext("123.123"), "{floor {0}}", "123")
	testExpression(t, mockContext("123.123"), "{ceil {0}}", "124")
	testExpression(t, mockContext("123.123"), "{round {0}}", "123")
	testExpression(t, mockContext("123.123"), "{round {0} 1}", "123.1")
	testExpression(t, mockContext("123.126"), "{round {0} 2}", "123.13")

	testExpressionErr(t, mockContext("123.123"), "{floor {0} b}", "<ARGN>", ErrArgCount)
	testExpressionErr(t, mockContext("123.123"), "{round {0} 1 2}", "<ARGN>", ErrArgCount)
	testExpressionErr(t, mockContext("123.123"), "{round {0} {0}}", "<CONST>", ErrConst)
	testExpressionErr(t, mockContext("123.123"), "{round {0} b}", "<CONST>", ErrConst)
}

func TestLogPow(t *testing.T) {
	testExpression(t, mockContext("100"), "{log10 {0}}", "2")
	testExpression(t, mockContext("64"), "{log2 {0}}", "6")
	testExpression(t, mockContext("64"), "{round {ln {0}} 4}", "4.1589")
	testExpression(t, mockContext("3"), "{pow {0} 3}", "27")
	testExpression(t, mockContext("81"), "{sqrt {0}}", "9")

	testExpressionErr(t, mockContext(), "{sqrt 1 2}", "<ARGN>", ErrArgCount)
}

// BenchmarkSumf/{sumf_0_1_2_3_4}-4         	 2699511	       429.5 ns/op	      26 B/op	       2 allocs/op

// 2 args
// Old:                BenchmarkSumf/{sumf_{0}_{0}}-4         	 3916035	       289.7 ns/op	      24 B/op	       1 allocs/op
// Dynamic non-static: BenchmarkSumf/{sumf_{0}_{0}}-4         	 3949767	       296.8 ns/op	      24 B/op	       1 allocs/op
// Dynamic single:     BenchmarkSumf/{sumf_{0}_1}-4         	 4475090	       261.7 ns/op	      24 B/op	       1 allocs/op
// Dynamic static:     BenchmarkSumf/{sumf_1_1}-4         	 5338790	       220.0 ns/op	      24 B/op	       1 allocs/op
func BenchmarkSumf(b *testing.B) {
	benchmarkExpression(b, mockContext("1"), "{sumf 1 1}", "2")
}

// Old      : BenchmarkSumi/{sumi_1_1_1_1}-4         	21108159	        60.30 ns/op	       0 B/op	       0 allocs/op
// 1 var    : BenchmarkSumi/{sumi_{0}_1_1_1}-4         	21317667	        47.90 ns/op	       0 B/op	       0 allocs/op
// All typed: BenchmarkSumi/{sumi_1_1_1_1}-4         	31323100	        32.60 ns/op	       0 B/op	       0 allocs/op
func BenchmarkSumi(b *testing.B) {
	benchmarkExpression(b, mockContext("1"), "{sumi 1 1 1 1}", "4")
}
