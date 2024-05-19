package expressions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTruthy(t *testing.T) {
	assert.True(t, Truthy(TruthyVal))
	assert.False(t, Truthy(FalsyVal))
	assert.False(t, Truthy("  "))
}

func TestTruthyStr(t *testing.T) {
	assert.Equal(t, TruthyVal, TruthyStr(true))
	assert.Equal(t, FalsyVal, TruthyStr(false))
}
