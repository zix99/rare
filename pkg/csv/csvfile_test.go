package csv

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleCSV(t *testing.T) {
	c, err := OpenCSV("-") // stdout
	assert.NoError(t, err)
	assert.NotNil(t, c)

	assert.NoError(t, c.Write([]string{"hello", "there"}))
	assert.NoError(t, c.WriteRow("hello", 1, int64(2)))

	assert.NoError(t, c.Close())
}

func TestCSVWrite(t *testing.T) {
	var buf bytes.Buffer
	c := NewCSV(&nopWriteCloser{&buf})
	c.WriteRow("hello", "thar")
	c.WriteRow("quack", "bob")
	c.Close()

	assert.Equal(t, "hello,thar\nquack,bob\n", buf.String())
}
