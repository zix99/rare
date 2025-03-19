package stdlib

import (
	"testing"
)

func TestCoalesce(t *testing.T) {
	testExpression(t,
		mockContext("", "a", "b"),
		"{coalesce {0}} {coalesce a b c} {coalesce {0} {2}}",
		" a b")
}

func TestBucketing(t *testing.T) {
	testExpression(t, mockContext("ab", "cd", "123"), "{bucket {2} 10} is bucketed", "120 is bucketed")
	testExpression(t, mockContext(), "{bucket -25 50}", "-50")
	testExpressionErr(t, mockContext(), "{bucket 70 -50}", "<VALUE>", ErrValue)
	testExpressionErr(t, mockContext(), "{bucket 5 a}", "<BAD-TYPE>", ErrNum)
	testExpressionErr(t, mockContext(), "{bucket 5}", "<ARGN>", ErrArgCount)
}

func TestBucketRange(t *testing.T) {
	testExpression(t, mockContext(), "{bucketrange 25 50}", "0 - 49")
	testExpression(t, mockContext(), "{bucketrange -25 50}", "-50 - -1")
	testExpression(t, mockContext(70), "{bucketrange {0} 50}", "50 - 99")
	testExpressionErr(t, mockContext(), "{bucketrange 70 -50}", "<VALUE>", ErrValue)
	testExpressionErr(t, mockContext(), "{bucketrange 5 a}", "<BAD-TYPE>", ErrNum)
	testExpressionErr(t, mockContext(), "{bucketrange 5}", "<ARGN>", ErrArgCount)
}

func TestBucket(t *testing.T) {
	testExpression(t,
		mockContext("1000", "1200", "1234"),
		"{bucket {0} 1000} {bucket {1} 1000} {bucket {2} 1000} {bucket {2} 100}",
		"1000 1000 1000 1200")
	testExpressionErr(t, mockContext(), "{bucket abc 100} {bucket 1}", "<BAD-TYPE> <ARGN>", ErrArgCount)
}

func TestExpBucket(t *testing.T) {
	testExpression(t, mockContext("123", "1234", "12345"),
		"{expbucket {0}} {expbucket {1}} {expbucket {2}}", "100 1000 10000")
}

func TestClamp(t *testing.T) {
	testExpression(t, mockContext("100", "200", "1000", "-10"),
		"{clamp {0} 50 200}-{clamp {1} 50 200}-{clamp {2} 50 200}-{clamp {3} 50 200}",
		"100-200-max-min")
	testExpressionErr(t, mockContext("0"), "{clamp {0} {0} 1}", "<BAD-TYPE>", ErrNum)
	testExpressionErr(t, mockContext("0"), "{clamp {0} 1 {0}}", "<BAD-TYPE>", ErrNum)
	testExpressionErr(t, mockContext("0"), "{clamp {0} 1 2 3}", "<ARGN>", ErrArgCount)
}
