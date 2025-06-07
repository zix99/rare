package expressions

import (
	"bytes"
	"errors"
	"strconv"
	"testing"
	"text/template"

	"github.com/zix99/rare/pkg/testutil"

	"github.com/stretchr/testify/assert"
)

var testContext = KeyBuilderContextArray{
	Elements: []string{"ab", "cd", "123"},
	Keys: map[string]string{
		"test": "testval",
	},
}

func TestSimpleKey(t *testing.T) {
	kb, err := NewKeyBuilder().Compile("test 123")
	key := kb.BuildKey(&testContext)
	assert.Nil(t, err)
	assert.Equal(t, "test 123", key)
	assert.Equal(t, 1, len(kb.stages))
}

func TestSimpleReplacement(t *testing.T) {
	kb, _ := NewKeyBuilder().Compile("{0} is {1}")
	key := kb.BuildKey(&testContext)
	assert.Equal(t, "ab is cd", key)
	assert.Equal(t, 3, len(kb.stages))
}

func TestEmpty(t *testing.T) {
	kb, _ := NewKeyBuilder().Compile("")
	key := kb.BuildKey(&testContext)

	assert.Zero(t, kb.StageCount())
	assert.Equal(t, "", key)
}

func TestUnterminatedReplacement(t *testing.T) {
	kb, err := NewKeyBuilder().Compile("{0} is {123")
	assert.Error(t, err)
	assert.Len(t, err.Errors, 1)
	assert.NotEmpty(t, err.Error())
	assert.NotNil(t, kb) // Still returns workable expression, but with errors
}

func TestManyErrors(t *testing.T) {
	kb, err := NewKeyBuilder().Compile("{0} is {abc 1} and {unclosed")
	assert.NotNil(t, kb)
	assert.Error(t, err)
	assert.Len(t, err.Errors, 2)

	// Test individually
	assert.ErrorIs(t, err.Errors[0], ErrorMissingFunction)
	assert.ErrorIs(t, err.Errors[1], ErrorUnterminated)

	// Test unwrap/is
	assert.ErrorIs(t, err, ErrorMissingFunction)
	assert.ErrorIs(t, err, ErrorUnterminated)
	assert.NotErrorIs(t, err, ErrorEmptyStatement)
	assert.ErrorIs(t, errors.Unwrap(err), ErrorMissingFunction)
	assert.NotEmpty(t, err.Error())
}

func TestEscapedString(t *testing.T) {
	kb, _ := NewKeyBuilder().Compile("{0} is \\{1\\} cool\\n\\t\\a\\r")
	key := kb.BuildKey(&testContext)
	assert.Equal(t, "ab is {1} cool\n\ta\r", key)
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

func TestEmptyStatement(t *testing.T) {
	kb, err := NewKeyBuilder().Compile("{} test")
	assert.NotNil(t, kb)
	assert.Error(t, err)
}

// BenchmarkSimpleReplacement-4   	 7515498	       141.4 ns/op	      24 B/op	       2 allocs/op
func BenchmarkSimpleReplacement(b *testing.B) {
	kb, _ := NewKeyBuilder().Compile("{0} is awesome")
	for n := 0; n < b.N; n++ {
		kb.BuildKey(&testContext)
	}
}

// BenchmarkGoTextTemplate-4   	 3139363	       406.3 ns/op	     160 B/op	       3 allocs/op
func BenchmarkGoTextTemplate(b *testing.B) {
	kb, _ := template.New("test").Parse("{a} is awesome")
	for n := 0; n < b.N; n++ {
		var buf bytes.Buffer
		kb.Execute(&buf, nil)
	}
}

// func tests

var simpleFuncs = map[string]KeyBuilderFunction{
	"addi": func(args []KeyBuilderStage) (KeyBuilderStage, error) {
		if len(args) < 2 {
			return nil, errors.New("expected at least 2 args")
		}
		return func(ctx KeyBuilderContext) string {
			val, _ := strconv.Atoi(args[0](ctx))
			for i := 1; i < len(args); i++ {
				aVal, _ := strconv.Atoi(args[i](ctx))
				val += aVal
			}
			return strconv.Itoa(val)
		}, nil
	},
}

func TestSimpleFuncs(t *testing.T) {
	k := NewKeyBuilderEx(false)
	k.Funcs(simpleFuncs)
	assert.True(t, k.HasFunc("addi"))
	assert.False(t, k.HasFunc("addb"))

	kb, _ := k.Compile("value: {addi {addi 1 2} 2}")
	assert.Equal(t, 2, kb.StageCount())
	assert.Equal(t, "value: 5", kb.BuildKey(&KeyBuilderContextArray{}))
}

func TestSimpleFuncErrors(t *testing.T) {
	k := NewKeyBuilder()
	k.Funcs(simpleFuncs)
	kb, err := k.Compile("value: {addi 1} {addi 1 2}")
	assert.Error(t, err)
	assert.NotNil(t, kb)
	assert.Equal(t, "value:  3", kb.BuildKey(&KeyBuilderContextArray{}))
}

func TestDeepFuncError(t *testing.T) {
	k := NewKeyBuilder()
	k.Funcs(simpleFuncs)
	kb, err := k.Compile("value: {addi 1 {addi 1}} {addi 1 2}")
	assert.Error(t, err)
	assert.NotNil(t, kb)
	assert.Equal(t, "value: 1 3", kb.BuildKey(&KeyBuilderContextArray{}))
}

func TestManyStages(t *testing.T) {
	k := NewKeyBuilderEx(false)
	k.Funcs(simpleFuncs)
	kb, _ := k.Compile("value: {addi -{addi 1 2} 2} {addi 3 5}")
	assert.Equal(t, 4, kb.StageCount())
	assert.Equal(t, "value: -1 8", kb.BuildKey(&KeyBuilderContextArray{}))
}

// Optimization

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

func TestOptimizedZeroAllocs(t *testing.T) {
	testutil.AssertZeroAlloc(t, BenchmarkOptimizedFunc)
}
