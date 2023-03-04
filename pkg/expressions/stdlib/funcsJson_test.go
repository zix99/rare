package stdlib

import (
	"testing"
)

func TestJson(t *testing.T) {
	testExpression(t, mockContext(`{"abc":123}`), `{json {0} abc}`, "123")
}

func TestJsonSingleArg(t *testing.T) {
	testExpression(t, mockContext(`{"abc":456}`), `{json abc}`, "456")
}

func TestJsonManyArgs(t *testing.T) {
	testExpression(t, mockContext(`{"abc":456}`), `{json {0} abc woops}`, "<ARGN>")
}

func TestJsonComplexObject(t *testing.T) {
	testExpression(t, mockContext(`{"abc":{"efg":23}}`), `{json {0} abc.efg}`, "23")
	testExpression(t, mockContext(`{"abc":{"efg":23}}`), `{json {0} abc.qef}`, "")
}

func TestJsonNestedArray(t *testing.T) {
	testExpression(t, mockContext(`[1,2,3,4]`), `{json 1}`, "2")
	testExpression(t, mockContext(`{"a":[1,2,3,4]}`), `{json a.1}`, "2")
	testExpression(t, mockContext(`{"a":[{"efg":123},2,3,4]}`), `{json a.0}`, `{"efg":123}`)
	testExpression(t, mockContext(`{"a":[{"efg":123},2,3,4]}`), `{json a.0.efg}`, `123`)
}

// BenchmarkJson-4   	 7041579	       169.1 ns/op	       0 B/op	       0 allocs/op
func BenchmarkJson(b *testing.B) {
	kb, _ := NewStdKeyBuilder().Compile("{json abc}")
	context := mockContext(`{"abc":123}`)
	for i := 0; i < b.N; i++ {
		kb.BuildKey(context)
	}
}
