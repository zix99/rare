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

func TestSimpleObjPoolCustomNew(t *testing.T) {
	type testObj struct{ item int }

	op := NewObjectPoolEx[testObj](1, func() *testObj {
		return &testObj{5}
	})

	assert.Equal(t, 5, op.Get().item)
	assert.Equal(t, 5, op.Get().item)
}

func TestZeroAllocs(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	res := testing.Benchmark(BenchmarkObjPool)
	assert.Zero(t, res.AllocsPerOp())
}

// BenchmarkObjPool-4   	24965836	        44.38 ns/op	       0 B/op	       0 allocs/op
func BenchmarkObjPool(b *testing.B) {
	type testObj struct{}

	op := NewObjectPool[testObj](5)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		op.Return(op.Get())
	}
}
