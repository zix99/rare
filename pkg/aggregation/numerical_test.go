package aggregation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleNumericalAggregation(t *testing.T) {
	aggr := NewNumericalAggregator(&NumericalConfig{})
	aggr.Samplef(5)
	aggr.Samplef(10)
	aggr.Samplef(15)

	assert.Equal(t, uint64(3), aggr.Count())
	assert.Equal(t, 10.0, aggr.Mean())
	assert.Equal(t, 5.0, aggr.Min())
	assert.Equal(t, 15.0, aggr.Max())
	assert.InEpsilon(t, 4.08248, aggr.StdDev(), 0.001)

	data := aggr.Analyze()

	assert.Equal(t, 10.0, data.Median())
	assert.Equal(t, 10.0, data.Quantile(0.5))
	assert.Equal(t, 5.0, data.Mode())
}

func TestSimpleMode(t *testing.T) {
	aggr := NewNumericalAggregator(&NumericalConfig{})
	aggr.Samplef(5)
	aggr.Samplef(10)
	aggr.Samplef(15)
	aggr.Samplef(5)
	aggr.Samplef(10)
	aggr.Samplef(5)

	data := aggr.Analyze()
	assert.Equal(t, 5.0, data.Mode())
	assert.Equal(t, 15.0, data.Quantile(0.9))
}
