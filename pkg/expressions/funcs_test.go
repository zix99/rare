package expressions

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testFuncData = []string{"ab", "100", "1000000", "5000000.123456", "22"}
var testFuncContext = KeyBuilderContextArray{
	Elements: testFuncData,
}

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
	kb, err := NewKeyBuilder().Compile(expression)
	assert.NoError(t, err)
	ret := kb.BuildKey(context)
	assert.Equal(t, expected, ret)
}

func TestCoalesce(t *testing.T) {
	testExpression(t,
		mockContext("", "a", "b"),
		"{coalesce a b c} {coalesce {0} {2}}",
		"a b")
}

func TestSimpleFunction(t *testing.T) {
	kb, _ := NewKeyBuilder().Compile("{hi {2}} {hf {3}}")
	key := kb.BuildKey(&testFuncContext)
	assert.Equal(t, "1,000,000 5,000,000.1235", key)
}

func TestByteSize(t *testing.T) {
	kb, _ := NewKeyBuilder().Compile("{bytesize {2}}")
	key := kb.BuildKey(&testFuncContext)
	assert.Equal(t, "976 KB", key)
}

func TestComparisonExpression(t *testing.T) {
	kb, _ := NewKeyBuilder().Compile("{and {lt {2} 10000000} {gt {1} 50}}")
	key := kb.BuildKey(&testFuncContext)
	assert.Equal(t, "1", key)
}

func TestNotExpression(t *testing.T) {
	kb, _ := NewKeyBuilder().Compile("{not {and {lt {2} 10000000} {gt {1} 50}}}")
	key := kb.BuildKey(&testFuncContext)
	assert.Equal(t, "", key)
}

func TestLike(t *testing.T) {
	kb, _ := NewKeyBuilder().Compile("{like {0} \"a\"}{like {0} c}")
	key := kb.BuildKey(&testFuncContext)
	assert.Equal(t, "ab", key)
}

func TestArithmatic(t *testing.T) {
	kb, _ := NewKeyBuilder().Compile("{sumi {1} {4}} {multi {1} 2} {divi {1} 2} {subi {1} 10}")
	key := kb.BuildKey(&testFuncContext)
	assert.Equal(t, "122 200 50 90", key)
}

func TestArithmaticf(t *testing.T) {
	kb, _ := NewKeyBuilder().Compile("{sumf {1} {4}} {multf {1} 2} {divf {1} 2} {subf {1} 10}")
	key := kb.BuildKey(&testFuncContext)
	assert.Equal(t, "122 200 50 90", key)
}
