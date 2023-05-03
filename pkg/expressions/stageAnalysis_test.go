package expressions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testStageUseContext(ret string) KeyBuilderStage {
	return func(context KeyBuilderContext) string {
		context.GetMatch(0)
		return ret
	}
}

func testStageNoContext(ret string) KeyBuilderStage {
	return func(context KeyBuilderContext) string {
		return ret
	}
}

func TestEvaluateStageIndex(t *testing.T) {
	stages := []KeyBuilderStage{
		testStageUseContext("test1"),
		testStageNoContext("test2"),
	}

	assert.Equal(t, "nope", EvalStageIndexOrDefault(stages, 0, "nope"))
	assert.Equal(t, "test2", EvalStageIndexOrDefault(stages, 1, "nope"))
	assert.Equal(t, "nope", EvalStageIndexOrDefault(stages, 2, "nope"))
}

func TestEvaluationStageInt(t *testing.T) {
	val, ok := EvalStageInt(testStageNoContext("5"))
	assert.Equal(t, 5, val)
	assert.True(t, ok)

	val, ok = EvalStageInt(testStageNoContext("5b"))
	assert.Equal(t, 0, val)
	assert.False(t, ok)

	val, ok = EvalStageInt(testStageUseContext("5"))
	assert.Equal(t, 0, val)
	assert.False(t, ok)
}

func TestEvaluationStageInt64(t *testing.T) {
	val, ok := EvalStageInt64(testStageNoContext("5"))
	assert.Equal(t, int64(5), val)
	assert.True(t, ok)

	val, ok = EvalStageInt64(testStageNoContext("5b"))
	assert.Equal(t, int64(0), val)
	assert.False(t, ok)

	val, ok = EvalStageInt64(testStageUseContext("5"))
	assert.Equal(t, int64(0), val)
	assert.False(t, ok)
}
