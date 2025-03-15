package stdlib

import (
	"rare/pkg/expressions"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypedStageStatic(t *testing.T) {
	s, ok := evalTypedStage(func(kbc expressions.KeyBuilderContext) string {
		return "1.0"
	}, typedParserFloat)
	assert.True(t, ok)

	v, vok := s(mockContext("2"))
	assert.True(t, vok)
	assert.Equal(t, 1.0, v)
}

func TestTypedStageStaticError(t *testing.T) {
	s, ok := evalTypedStage(func(kbc expressions.KeyBuilderContext) string {
		return "blabla"
	}, typedParserFloat)
	assert.False(t, ok)
	assert.Nil(t, s)
}

func TestTypedStageDynamic(t *testing.T) {
	s, ok := evalTypedStage(func(kbc expressions.KeyBuilderContext) string {
		return kbc.GetMatch(0)
	}, typedParserFloat)
	assert.True(t, ok)

	v, vok := s(mockContext("1.0"))
	assert.True(t, vok)
	assert.Equal(t, 1.0, v)
}

func TestTypedStageDynamicError(t *testing.T) {
	s, ok := evalTypedStage(func(kbc expressions.KeyBuilderContext) string {
		return kbc.GetMatch(0)
	}, typedParserFloat)
	assert.True(t, ok)

	_, vok := s(mockContext("abc"))
	assert.False(t, vok)
}

func TestMapTypedStages(t *testing.T) {
	stages := []expressions.KeyBuilderStage{
		func(kbc expressions.KeyBuilderContext) string {
			return "1"
		},
		func(kbc expressions.KeyBuilderContext) string {
			return kbc.GetMatch(0)
		},
		func(kbc expressions.KeyBuilderContext) string {
			return kbc.GetMatch(1)
		},
	}

	mstages, ok := mapTypedArgs(stages, typedParsedInt)
	assert.True(t, ok)

	ctx := mockContext("5", "bla")

	arg0, ok0 := mstages[0](ctx)
	assert.True(t, ok0)
	assert.Equal(t, 1, arg0)

	arg1, ok1 := mstages[1](ctx)
	assert.True(t, ok1)
	assert.Equal(t, 5, arg1)

	arg2, ok2 := mstages[2](ctx)
	assert.False(t, ok2)
	assert.Equal(t, 0, arg2)
}
