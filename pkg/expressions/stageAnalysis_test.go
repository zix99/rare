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

func TestEvaluateStageStatic(t *testing.T) {
	stage := testStageNoContext("test")
	assert.Equal(t, "test", EvalStageOrDefault(stage, "nope"))
}

func TestEvaluateStageDynamic(t *testing.T) {
	stage := testStageUseContext("test")
	assert.Equal(t, "nope", EvalStageOrDefault(stage, "nope"))
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
	val, ok := EvalStageInt(testStageNoContext("5"), 1)
	assert.Equal(t, 5, val)
	assert.True(t, ok)

	val, ok = EvalStageInt(testStageNoContext("5b"), 1)
	assert.Equal(t, 1, val)
	assert.False(t, ok)

	val, ok = EvalStageInt(testStageUseContext("5"), 1)
	assert.Equal(t, 1, val)
	assert.False(t, ok)
}
