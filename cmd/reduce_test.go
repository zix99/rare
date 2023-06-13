package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReduce(t *testing.T) {
	out, eout, err := testCommandCapture(reduceCommand(),
		`-m (\d+) --snapshot -a "test={sumi {.} {0}}" testdata/log.txt`)
	assert.NoError(t, err)
	assert.Empty(t, eout)
	assert.Equal(t, "test: 32\nMatched: 3 / 6\n96 B (0 B/s) \n", out)
}

func TestReduceBasics(t *testing.T) {
	testCommandSet(t, reduceCommand(),
		`-m (\d+) -g {0} -a {0} testdata/log.txt`,
		`-m (\d+) -g {0} -a a={0} --sort {a} --sort-reverse testdata/log.txt`,
		`-o - -m (\d+) -g {0} -a a={0} --sort {a} --sort-reverse testdata/log.txt`)
}

func TestParseKIV(t *testing.T) {
	k, i, v := parseKeyValInitial("abc", "init")
	assert.Equal(t, "abc", k)
	assert.Equal(t, "init", i)
	assert.Equal(t, "abc", v)

	k, i, v = parseKeyValInitial("abc=efg", "init")
	assert.Equal(t, "abc", k)
	assert.Equal(t, "init", i)
	assert.Equal(t, "efg", v)

	k, i, v = parseKeyValInitial("=efg", "init")
	assert.Equal(t, "", k)
	assert.Equal(t, "init", i)
	assert.Equal(t, "efg", v)

	k, i, v = parseKeyValInitial("abc:=efg", "init")
	assert.Equal(t, "abc", k)
	assert.Equal(t, "", i)
	assert.Equal(t, "efg", v)

	k, i, v = parseKeyValInitial("abc:1=efg", "init")
	assert.Equal(t, "abc", k)
	assert.Equal(t, "1", i)
	assert.Equal(t, "efg", v)
}
