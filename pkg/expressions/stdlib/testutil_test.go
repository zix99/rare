package stdlib

import (
	"fmt"
	. "rare/pkg/expressions"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockContext(args ...interface{}) KeyBuilderContext {
	s := make([]string, len(args))
	for idx, arg := range args {
		s[idx] = fmt.Sprintf("%v", arg)
	}
	return &KeyBuilderContextArray{
		Elements: s,
	}
}

func testExpression(t *testing.T, context KeyBuilderContext, expression string, expected string) {
	kb, err := NewStdKeyBuilder().Compile(expression)
	assert.Nil(t, err)
	assert.NotNil(t, kb)
	if kb != nil {
		ret := kb.BuildKey(context)
		assert.Equal(t, expected, ret)
	}
}

// if expected is nil, any error is acceptable
func testExpressionErr(t *testing.T, context KeyBuilderContext, expression string, evalsTo string, expected ...interface{}) {
	kb, err := NewStdKeyBuilder().Compile(expression)
	if evalsTo == "" {
		assert.Nil(t, kb, "Expected not compiled with error")
	} else {
		assert.NotNil(t, kb, "Expected compiled with error")
		if kb != nil {
			assert.Equal(t, evalsTo, kb.BuildKey(context))
		}
	}

	assert.NotNil(t, err, "Expected error")
	if len(expected) > 0 && err != nil {
		for _, ex := range expected {
			switch e := ex.(type) {
			case funcError:
				assert.ErrorIs(t, err, e.err)
			case error:
				assert.ErrorIs(t, err, e)
			default:
				t.Error("Invalid type assertion, expected error or funcError")
			}
		}
	}
}

// benchmark an expreession, as a sub-benchmark. Checks value before running test
func benchmarkExpression(b *testing.B, context KeyBuilderContext, expression, expected string) {
	kb, err := NewStdKeyBuilderEx(false).Compile(expression)
	if err != nil {
		b.Fatal(err)
	}

	if s := kb.BuildKey(context); s != expected {
		b.Fatalf("%s != %s", s, expected)
	}

	b.Run(expression, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			kb.BuildKey(context)
		}
	})
}
