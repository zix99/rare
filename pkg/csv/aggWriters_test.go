package csv

import (
	"bytes"
	"rare/pkg/aggregation"
	"rare/pkg/expressions"
	"rare/pkg/expressions/stdlib"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteTable(t *testing.T) {
	table := aggregation.NewTable(" ")
	table.Sample("hello there")
	table.Sample("hello there")
	table.Sample("hello bob")
	table.Sample("cat bob")

	var buf bytes.Buffer
	csv := NewCSV(&nopWriteCloser{&buf})
	assert.NoError(t, WriteTable(csv, table))
	csv.Close()

	assert.Equal(t, ",cat,hello\nbob,1,1\nthere,0,2\n", buf.String())
}

func TestWriteAccumulator(t *testing.T) {
	accum := aggregation.NewAccumulatingGroup(stdlib.NewStdKeyBuilder())
	accum.AddGroupExpr("g0", "{1}")
	accum.AddDataExpr("d0", "{1}", "0")
	accum.AddDataExpr("d1", "{sumi {.} {3}}", "0")

	accum.Sample(expressions.MakeArray("1", "2", "3"))
	accum.Sample(expressions.MakeArray("1", "2", "3"))
	accum.Sample(expressions.MakeArray("2", "2", "3"))

	var buf bytes.Buffer
	csv := NewCSV(&nopWriteCloser{&buf})
	assert.NoError(t, WriteAccumulator(csv, accum))
	csv.Close()

	assert.Equal(t, "g0,d0,d1\n1,1,6\n2,2,3\n", buf.String())
}

func TestWriteCounter(t *testing.T) {
	counter := aggregation.NewCounter()
	counter.Sample("bla")
	counter.Sample("bla")
	counter.Sample("cookie")

	var buf bytes.Buffer
	csv := NewCSV(&nopWriteCloser{&buf})
	assert.NoError(t, WriteCounter(csv, counter))
	csv.Close()

	assert.Equal(t, "group,value\nbla,2\ncookie,1\n", buf.String())
}

func TestWriteSubCounter(t *testing.T) {
	counter := aggregation.NewSubKeyCounter()
	counter.Sample(expressions.MakeArray("hello", "there"))
	counter.Sample(expressions.MakeArray("hello", "there"))
	counter.Sample(expressions.MakeArray("hello", "bob"))
	counter.Sample(expressions.MakeArray("bob", "hello"))

	var buf bytes.Buffer
	csv := NewCSV(&nopWriteCloser{&buf})
	assert.NoError(t, WriteSubCounter(csv, counter))
	csv.Close()

	assert.Equal(t, "group,bob,hello,there\nbob,0,1,0\nhello,1,0,2\n", buf.String())
}
