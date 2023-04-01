package aggregation

import (
	"rare/pkg/expressions"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleAccumulator(t *testing.T) {
	accum, err := NewExprAccumulator("{sumi {.} 1}", "5")
	assert.NoError(t, err)
	assert.Equal(t, "5", accum.Value())

	accum.Sample("hello there")
	assert.Equal(t, "6", accum.Value())
}

func TestAccumulateValues(t *testing.T) {
	accum, err := NewExprAccumulator("{sumi {.} {0}}", "0")
	assert.NoError(t, err)
	assert.Equal(t, "0", accum.Value())

	accum.Sample("2")
	accum.Sample("1")
	assert.Equal(t, "3", accum.Value())
}

func TestAccumulateRange(t *testing.T) {
	accum, err := NewExprAccumulator("{sumi {.} {1} {2}}", "0")
	assert.NoError(t, err)
	assert.Equal(t, "0", accum.Value())

	accum.Sample(expressions.MakeArray("1", "2"))
	accum.Sample(expressions.MakeArray("3", "4"))
	assert.Equal(t, "10", accum.Value())
}

func TestAccumulateSet(t *testing.T) {
	aset := NewExprAccumulatorSet()
	assert.NoError(t, aset.Add("sum", "{sumi {.} {1}}", "0"))
	assert.NoError(t, aset.Add("mult", "{multi {.} {1}}", "1"))
	assert.NoError(t, aset.Add("resum", "{sumi {sum} {mult} 1}", "0"))

	aset.Sample("2")
	aset.Sample("2")
	aset.Sample("3")

	items := aset.Items()
	assert.Equal(t, "sum", items[0].Name)
	assert.Equal(t, "7", items[0].Accum.Value())

	assert.Equal(t, "mult", items[1].Name)
	assert.Equal(t, "12", items[1].Accum.Value())

	assert.Equal(t, "resum", items[2].Name)
	assert.Equal(t, "20", items[2].Accum.Value())
}
