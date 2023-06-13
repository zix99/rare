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
	assert.Equal(t, accum.DataCount(), 1)
	assert.Equal(t, "4", accum.Data("")[0])
	assert.Equal(t, "4", accum.DataNoCopy("")[0])
	assert.Zero(t, accum.ParseErrors())
}

func TestAccumGroups(t *testing.T) {
	accum := NewAccumulatingGroup(stdlib.NewStdKeyBuilder())

	accum.AddGroupExpr("test", "{1}")
	accum.AddDataExpr("sum", "{sumi {.} {2}}", "0")
	accum.AddDataExpr("mul", "{multi {.} {2}}", "1")

	accum.Sample(expressions.MakeArray("200", "2"))
	accum.Sample(expressions.MakeArray("200", "3"))
	accum.Sample(expressions.MakeArray("400", "2"))

	assert.Len(t, accum.GroupCols(), 1)
	assert.Equal(t, 1, accum.GroupColCount())
	assert.Equal(t, 2, len(accum.data))
	assert.Equal(t, 2, accum.DataCount())
}

func TestMultiGroupCols(t *testing.T) {
	accum := NewAccumulatingGroup(stdlib.NewStdKeyBuilder())

	accum.AddGroupExpr("test", "{1}")
	accum.AddGroupExpr("test2", "{bucket {2} 10}")
	accum.AddDataExpr("sum", "{sumi {.} {2}}", "0")
	accum.AddDataExpr("mul", "{multi {.} {2}}", "1")

	accum.Sample(expressions.MakeArray("200", "2"))
	accum.Sample(expressions.MakeArray("200", "3"))
	accum.Sample(expressions.MakeArray("400", "2"))

	assert.Len(t, accum.GroupCols(), 2)
	assert.Equal(t, 2, accum.GroupColCount())
	assert.Equal(t, 4, accum.ColCount())
	assert.Equal(t, 2, len(accum.data))
	assert.Equal(t, []GroupKey{"200\x000", "400\x000"}, accum.Groups(sorting.ByName))
}

func TestFaltiGroupMatch(t *testing.T) {
	accum := NewAccumulatingGroup(stdlib.NewStdKeyBuilder())
	assert.NoError(t, accum.AddGroupExpr("test", "{badkey}"))
	accum.Sample("100")
	assert.Equal(t, 1, accum.GroupColCount())
	assert.Equal(t, []GroupKey{""}, accum.Groups(sorting.ByName))
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

	assert.NoError(t, accum.AddDataExpr("test", "{sumi {.} {bla}}", "0"))
	assert.Error(t, accum.AddDataExpr("test", "{0}", "0")) // Dupe key error

	assert.NoError(t, accum.AddGroupExpr("dupe", "{0}"))
	assert.Error(t, accum.AddGroupExpr("dupe", "{0}")) // Dupe group

	// Sample
	accum.Sample("123")
	assert.Equal(t, accum.Data("123")[0], "<BAD-TYPE>")

	assert.Error(t, accum.AddDataExpr("real", "{0}", "0"))
	assert.Error(t, accum.AddGroupExpr("real", "{0}"))

}

func TestAccumSort(t *testing.T) {
	accum := NewAccumulatingGroup(stdlib.NewStdKeyBuilder())

	accum.AddGroupExpr("test", "{1}")
	accum.AddDataExpr("sum", "{sumi {.} {2}}", "0")
	accum.AddDataExpr("mul", "{multi {.} {2}}", "1")

	accum.Sample(expressions.MakeArray("200", "2"))
	accum.Sample(expressions.MakeArray("200", "3"))
	accum.Sample(expressions.MakeArray("400", "2"))
	accum.Sample(expressions.MakeArray("800", "1"))

	assert.NoError(t, accum.SetSort("{sum}"))
	assert.Equal(t, []GroupKey{"800", "400", "200"}, accum.Groups(sorting.ByNameSmart))

	assert.NoError(t, accum.SetSort("{.}"))
	assert.Equal(t, []GroupKey{"200", "400", "800"}, accum.Groups(sorting.ByNameSmart))

	assert.NoError(t, accum.SetSort("-{0}"))
	assert.Equal(t, []GroupKey{"800", "400", "200"}, accum.Groups(sorting.ByNameSmart))

	assert.Error(t, accum.SetSort("{0"))
}

func BenchmarkAccumulatorContext(b *testing.B) {
	ctx := exprAccumulatorContext{
		current: "123",
		match:   "1\x002\x003",
	}
	for i := 0; i < b.N; i++ {
		ctx.GetMatch(3)
	}
}

// BenchmarkGroupKey-4   	33043279	        32.26 ns/op	       0 B/op	       0 allocs/op
func BenchmarkGroupKey(b *testing.B) {
	accum := NewAccumulatingGroup(stdlib.NewStdKeyBuilder())
	accum.AddGroupExpr("key", "{1}")

	ctx := exprAccumulatorContext{
		match: "123",
	}

	for i := 0; i < b.N; i++ {
		accum.buildGroupKey("hello\x00thar", &ctx)
	}
}

// BenchmarkSample-4   	 4278193	       278.8 ns/op	      51 B/op	       2 allocs/op
func BenchmarkSample(b *testing.B) {
	accum := NewAccumulatingGroup(stdlib.NewStdKeyBuilder())
	accum.AddGroupExpr("key", "{0}")
	accum.AddDataExpr("max", "{maxi {.} {0}}", "0")

	for i := 0; i < b.N; i++ {
		accum.Sample("123")
	}
}
