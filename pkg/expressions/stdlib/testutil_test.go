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
	kb, _ := NewStdKeyBuilder().Compile(expression)
	//assert.NoError(t, err)
	assert.NotNil(t, kb)
	if kb != nil {
		ret := kb.BuildKey(context)
		assert.Equal(t, expected, ret)
	}
}

// if expected is nil, any error is acceptable
func testExpressionErr(t *testing.T, context KeyBuilderContext, expression string, evalsTo string, expected ...error) {
	kb, err := NewStdKeyBuilder().Compile(expression)
	if evalsTo == "" {
		assert.Nil(t, kb, "Expected not compiled with error")
	} else {
		assert.NotNil(t, kb, "Expected compiled with error")
		if kb != nil {
			assert.Equal(t, evalsTo, kb.BuildKey(context))
		}
	}
	if len(expected) == 0 {
		assert.Error(t, err)
	} else {
		assert.ErrorIs(t, err, expected[0])
	}
}
