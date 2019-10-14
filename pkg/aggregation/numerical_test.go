package aggregation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleNumericalAggregation(t *testing.T) {
	aggr := NewNumericalAggregator()
	aggr.Sample(5)
	aggr.Sample(10)
	aggr.Sample(15)

	assert.Equal(t, uint64(3), aggr.Count())
	assert.Equal(t, 10.0, aggr.Mean())
	assert.Equal(t, 5.0, aggr.Min())
	assert.Equal(t, 15.0, aggr.Max())

	data := aggr.Analyze()

	assert.Equal(t, 10.0, data.Mean)
	assert.Equal(t, 10.0, data.Median)
}
