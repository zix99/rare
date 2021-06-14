package minijson

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyObject(t *testing.T) {
	var jb JsonObjectBuilder
	assert.Empty(t, jb.String())

	jb.Open()
	jb.Close()
	assert.Equal(t, "{}", jb.String())
}

func TestSingleObject(t *testing.T) {
	var jb JsonObjectBuilder
	jb.OpenEx(5)
	jb.WriteString("key", "val")
	jb.Close()
	assert.Equal(t, `{"key": "val"}`, jb.String())
}

func TestMultiObject(t *testing.T) {
	var jb JsonObjectBuilder
	jb.OpenEx(5)
	jb.WriteString("key", "val")
	jb.WriteString("key2", "val2")
	jb.Close()
	assert.Equal(t, `{"key": "val", "key2": "val2"}`, jb.String())
}

func TestWriteInt(t *testing.T) {
	var jb JsonObjectBuilder
	jb.Open()
	jb.WriteInt("k", 123)
	jb.Close()
	assert.Equal(t, `{"k": 123}`, jb.String())
}

func TestWriteInferred(t *testing.T) {
	var jb JsonObjectBuilder
	jb.Open()
	jb.WriteInferred("s", "\nstring")
	jb.WriteInferred("n", "123")
	jb.WriteInferred("f", "123.2")
	jb.WriteInferred("t", "True")
	jb.WriteInferred("fa", "FALSE")
	jb.Close()

	assert.Equal(t, `{"s": "\nstring", "n": 123, "f": 123.2, "t": true, "fa": false}`, jb.String())
}

func TestIsNumeric(t *testing.T) {
	assert.True(t, isNumeric("1"))
	assert.True(t, isNumeric("123"))
	assert.True(t, isNumeric("1.2"))
	assert.True(t, isNumeric("1000"))
	assert.False(t, isNumeric("0123a"))
	assert.False(t, isNumeric("qq"))
}

func TestEscape(t *testing.T) {
	assert.Equal(t, "this is normal", escape("this is normal"))
	assert.Equal(t, `\r\n new\t\t`, escape("\r\n new\t\t"))
}

var mapData = map[string]string{
	"s":  "string",
	"n":  "123",
	"f":  "123.2",
	"t":  "True",
	"fa": "FALSE",
}

func BenchmarkJsonMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		json.Marshal(mapData)
	}
}

// BenchmarkJsonBuilder-4   	  807200	      1344 ns/op	     256 B/op	       1 allocs/op
func BenchmarkJsonBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MarshalStringMapInferred(mapData)
	}
}

// BenchmarkEscape-4   	 1938726	       567.6 ns/op	      32 B/op	       1 allocs/op
func BenchmarkEscape(b *testing.B) {
	for i := 0; i < b.N; i++ {
		escape("\nthis is a test!")
	}
}
