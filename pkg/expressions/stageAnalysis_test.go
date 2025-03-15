package expressions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvalStaticHit(t *testing.T) {
	val, ok := EvalStaticStage(func(kbc KeyBuilderContext) string {
		return "static"
	})

	assert.True(t, ok)
	assert.Equal(t, "static", val)
}

func TestEvalStaticMissMatch(t *testing.T) {
	val, ok := EvalStaticStage(func(kbc KeyBuilderContext) string {
		return kbc.GetMatch(0) + "hit"
	})

	assert.False(t, ok)
	assert.Equal(t, "hit", val)
}

func TestEvalStaticMissKey(t *testing.T) {
	val, ok := EvalStaticStage(func(kbc KeyBuilderContext) string {
		return kbc.GetKey("stuff") + "hit"
	})

	assert.False(t, ok)
	assert.Equal(t, "hit", val)
}
