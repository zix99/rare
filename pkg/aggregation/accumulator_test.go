package aggregation

import (
	"rare/pkg/aggregation/sorting"
	"rare/pkg/expressions"
	"rare/pkg/expressions/stdlib"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyAccum(t *testing.T) {
	accum := NewAccumulatingGroup(stdlib.NewStdKeyBuilder())
	accum.Sample("1")
}

func TestBasicAccum(t *testing.T) {
	accum := NewAccumulatingGroup(stdlib.NewStdKeyBuilder())
	accum.AddDataExpr("sum", "{sumi {.} {0}}", "0")
	accum.Sample("1")
	accum.Sample("3")

	assert.Len(t, accum.GroupCols(), 0)
	assert.Len(t, accum.Groups(sorting.ByName), 1)
	assert.Len(t, accum.DataCols(), 1)
	assert.Equal(t, "4", accum.Data("")[0])
}

func TestAccumGroups(t *testing.T) {
	accum := NewAccumulatingGroup(stdlib.NewStdKeyBuilder())

	accum.AddGroupExpr("test", "{1}")
	accum.AddDataExpr("sum", "{sumi {.} {2}}", "0")
	accum.AddDataExpr("mul", "{multi {.} {2}}", "1")

	accum.Sample(expressions.MakeArray("200", "2"))
	accum.Sample(expressions.MakeArray("200", "3"))
	accum.Sample(expressions.MakeArray("400", "2"))

	assert.Equal(t, 1, accum.GroupColCount())
	assert.Equal(t, 2, len(accum.data))
}

func TestAccumSelfReference(t *testing.T) {
	accum := NewAccumulatingGroup(stdlib.NewStdKeyBuilder())

	accum.AddDataExpr("sum", "{sumi {.} {0}}", "0")
	accum.AddDataExpr("count", "{sumi {.} 1}", "0")
	accum.AddDataExpr("avg", "{divf {sum} {count}}", "")

	accum.Sample("4")
	accum.Sample("6")
	accum.Sample("10")
	accum.Sample("20")

	data := accum.Data("")
	assert.Equal(t, "40", data[0])
	assert.Equal(t, "4", data[1])
	assert.Equal(t, "10", data[2])
}

func TestParseGroupKey(t *testing.T) {
	assert.Equal(t, []string{}, GroupKey("").Parts())
	assert.Equal(t, []string{"a"}, GroupKey("a").Parts())
	assert.Equal(t, []string{"b", "c"}, GroupKey("b\x00c").Parts())
}

func TestAccumErrorCases(t *testing.T) {
	accum := NewAccumulatingGroup(stdlib.NewStdKeyBuilder())

	assert.Error(t, accum.AddDataExpr("", "{badexpr", ""))
	assert.Error(t, accum.AddGroupExpr("", "{badexpr"))

	accum.AddDataExpr("test", "{sumi {.} {bla}}", "0")
	accum.Sample("123")
	assert.Equal(t, accum.Data("")[0], "<BAD-TYPE>")
}
