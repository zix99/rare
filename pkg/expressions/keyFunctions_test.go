package expressions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testFuncData = []string{"ab", "100", "1000000", "5000000.123456", "22"}
var testFuncContext = KeyBuilderContextArray{
	Elements: testFuncData,
}

func TestSimpleFunction(t *testing.T) {
	kb := NewKeyBuilder().Compile("{hi {2}} {hf {3}}")
	key := kb.BuildKey(&testFuncContext)
	assert.Equal(t, "1,000,000 5,000,000.1235", key)
}

func TestByteSize(t *testing.T) {
	kb := NewKeyBuilder().Compile("{bytesize {2}}")
	key := kb.BuildKey(&testFuncContext)
	assert.Equal(t, "976 KB", key)
}

func TestExpression(t *testing.T) {
	kb := NewKeyBuilder().Compile("{and {lt {2} 10000000} {gt {1} 50}}")
	key := kb.BuildKey(&testFuncContext)
	assert.Equal(t, "1", key)
}

func TestNotExpression(t *testing.T) {
	kb := NewKeyBuilder().Compile("{not {and {lt {2} 10000000} {gt {1} 50}}}")
	key := kb.BuildKey(&testFuncContext)
	assert.Equal(t, "", key)
}

func TestLike(t *testing.T) {
	kb := NewKeyBuilder().Compile("{like {0} \"a\"}{like {0} c}")
	key := kb.BuildKey(&testFuncContext)
	assert.Equal(t, "ab", key)
}

func TestArithmatic(t *testing.T) {
	kb := NewKeyBuilder().Compile("{sumi {1} {4}} {multi {1} 2} {divi {1} 2} {subi {1} 10}")
	key := kb.BuildKey(&testFuncContext)
	assert.Equal(t, "122 200 50 90", key)
}
