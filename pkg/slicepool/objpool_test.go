package slicepool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleObjPool(t *testing.T) {
	type testObj struct{}

	op := NewObjectPool[testObj](1)
	assert.Len(t, op.pool, 1)
	v1 := op.Get()
	v2 := op.Get()
	assert.Len(t, op.pool, 0)
	assert.NotNil(t, v1)
	assert.NotNil(t, v2)

	op.Return(v1)
	op.Return(v2)
	assert.Len(t, op.pool, 2)
}
