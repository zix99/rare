package expressions

import (
	"bytes"
	"strconv"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

var testData = []string{"ab", "cd", "123"}
var testKeyData = map[string]string{
	"test": "testval",
}

type TestContext struct{}

func (s *TestContext) GetMatch(idx int) string {
	return testData[idx]
}

func (s *TestContext) GetKey(key string) string {
	return testKeyData[key]
}

var testContext = TestContext{}

func TestSimpleKey(t *testing.T) {
	kb, _ := NewKeyBuilder().Compile("test 123")
	key := kb.BuildKey(&testContext)
	assert.Equal(t, "test 123", key)
	assert.Equal(t, 1, len(kb.stages))
}

func TestSimpleReplacement(t *testing.T) {
	kb, _ := NewKeyBuilder().Compile("{0} is {1}")
	key := kb.BuildKey(&testContext)
	assert.Equal(t, "ab is cd", key)
	assert.Equal(t, 3, len(kb.stages))
}

func TestUnterminatedReplacement(t *testing.T) {
	kb, err := NewKeyBuilder().Compile("{0} is {123")
	assert.Error(t, err)
	assert.Nil(t, kb)
}

func TestEscapedString(t *testing.T) {
	kb, _ := NewKeyBuilder().Compile("{0} is \\{1\\} cool\\n\\t\\a")
	key := kb.BuildKey(&testContext)
	assert.Equal(t, "ab is {1} cool\n\ta", key)
	assert.Equal(t, 2, len(kb.stages))
}

func TestDeepKeys(t *testing.T) {
	kb, _ := NewKeyBuilder().Compile("{{1} b} is bucketed")
	key := kb.BuildKey(&testContext)
	assert.Equal(t, "<Err:{1}> is bucketed", key)
}

func TestStringKey(t *testing.T) {
	kb, _ := NewKeyBuilder().Compile("{test} {some} key")
	key := kb.BuildKey(&testContext)
	assert.Equal(t, "testval  key", key)
}

func BenchmarkSimpleReplacement(b *testing.B) {
	kb, _ := NewKeyBuilder().Compile("{0} is awesome")
	for n := 0; n < b.N; n++ {
		kb.BuildKey(&testContext)
	}
}

func BenchmarkGoTextTemplate(b *testing.B) {
	kb, _ := template.New("test").Parse("{a} is awesome")
	for n := 0; n < b.N; n++ {
		var buf bytes.Buffer
		kb.Execute(&buf, nil)
	}
}

// func tests

var simpleFuncs = map[string]KeyBuilderFunction{
	"addi": func(args []KeyBuilderStage) KeyBuilderStage {
		return func(ctx KeyBuilderContext) string {
			val, _ := strconv.Atoi(args[0](ctx))
			for i := 1; i < len(args); i++ {
				aVal, _ := strconv.Atoi(args[i](ctx))
				val += aVal
			}
			return strconv.Itoa(val)
		}
	},
}

func TestSimpleFuncs(t *testing.T) {
	k := NewKeyBuilderEx(false)
	k.Funcs(simpleFuncs)
	kb, _ := k.Compile("value: {addi {addi 1 2} 2}")
	assert.Equal(t, 2, kb.StageCount())
	assert.Equal(t, "value: 5", kb.BuildKey(&KeyBuilderContextArray{}))
}

func TestManyStages(t *testing.T) {
	k := NewKeyBuilderEx(false)
	k.Funcs(simpleFuncs)
	kb, _ := k.Compile("value: {addi -{addi 1 2} 2} {addi 3 5}")
	assert.Equal(t, 4, kb.StageCount())
	assert.Equal(t, "value: -1 8", kb.BuildKey(&KeyBuilderContextArray{}))
}

func TestManyStagesOptimize(t *testing.T) {
	k := NewKeyBuilderEx(true)

	k.Funcs(simpleFuncs)
	kb, _ := k.Compile("value: {addi -{addi 1 2} 2} {addi 3 5}")
	assert.Equal(t, 1, kb.StageCount())
	assert.Equal(t, "value: -1 8", kb.BuildKey(&KeyBuilderContextArray{}))
}

// BenchmarkSimpleFunc-4   	 9467798	       121.2 ns/op	       8 B/op	       1 allocs/op
func BenchmarkSimpleFunc(b *testing.B) {
	k := NewKeyBuilderEx(false)
	k.Funcs(simpleFuncs)
	kb, _ := k.Compile("value: {addi {addi 1 2} 2}")
	ctx := &KeyBuilderContextArray{}
	for i := 0; i < b.N; i++ {
		kb.BuildKey(ctx)
	}
}

// BenchmarkOptimizedFunc-4   	206643440	         5.698 ns/op	       0 B/op	       0 allocs/op
func BenchmarkOptimizedFunc(b *testing.B) {
	k := NewKeyBuilderEx(true)
	k.Funcs(simpleFuncs)
	kb, _ := k.Compile("value: {addi {addi 1 2} 2}")
	ctx := &KeyBuilderContextArray{}
	for i := 0; i < b.N; i++ {
		kb.BuildKey(ctx)
	}
}
