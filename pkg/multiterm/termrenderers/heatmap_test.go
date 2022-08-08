package termrenderers

import (
	"rare/pkg/aggregation"
	"rare/pkg/multiterm"
	"testing"
)

func TestSimpleHeatmap(t *testing.T) {
	hm := NewHeatmap(multiterm.NewVirtualTerm(), 10, 10)

	agg := aggregation.NewTable(" ")
	agg.Sample("test abc")

	hm.WriteTable(agg)

}
