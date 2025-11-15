package cmd

import (
	"testing"

	"github.com/zix99/rare/pkg/aggregation"
	"github.com/zix99/rare/pkg/expressions/funclib"
	"github.com/zix99/rare/pkg/logger"
	"github.com/zix99/rare/pkg/testutil"

	"github.com/stretchr/testify/assert"
)

func TestReduce(t *testing.T) {
	out, eout, err := testCommandCapture(reduceCommand(),
		`-m (\d+) --snapshot -a "test={sumi {.} {0}}" testdata/log.txt`)
	assert.NoError(t, err)
	assert.Empty(t, eout)
	testutil.AssertPattern(t, out, "test: 32\nMatched: 3 / 6\n96 B in * (*)\n")
}

func TestReduceBasics(t *testing.T) {
	testCommandSet(t, reduceCommand(),
		`-m (\d+) -g {0} -a {0} testdata/log.txt`,
		`-m (\d+) -g {0} -a a={0} --sort {a} --sort-reverse testdata/log.txt`,
		`-m (\d+) -g {0} -a a={0} --sort {a} --fmt downscale --sort-reverse testdata/log.txt`,
		`-m (\d+) -g {0} -a a={0} --sort {a} --fmt a=downscale --sort-reverse testdata/log.txt`,
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

func TestReduceFatals(t *testing.T) {
	catchLogFatal(t, 2, func() {
		testCommand(reduceCommand(), `-m (\d+) -g {0 -a {0} testdata/log.txt`)
	})
	catchLogFatal(t, 2, func() {
		testCommand(reduceCommand(), `-m (\d+) -g {0} -a {0 testdata/log.txt`)
	})
	catchLogFatal(t, 2, func() {
		testCommand(reduceCommand(), `-m (\d+) -g {0} -a {0} --sort {0 testdata/log.txt`)
	})
}

func TestBuildFormatterSet(t *testing.T) {
	accum := aggregation.NewAccumulatingGroup(funclib.NewKeyBuilder())

	accum.AddGroupExpr("by0", "{0}")
	accum.AddDataExpr("sum", "{sumi {.} {1}}", "0")
	accum.AddDataExpr("mult", "{multi {.} {1}}", "1")

	t.Run("default", func(t *testing.T) {
		deflt := buildFormatterSetOrFail(accum)
		assert.Len(t, deflt, 2)
	})

	t.Run("global", func(t *testing.T) {
		f := buildFormatterSetOrFail(accum, "bytesize")
		assert.Len(t, f, 2)
	})

	t.Run("byname", func(t *testing.T) {
		f := buildFormatterSetOrFail(accum, "bytesize", "sum=hi", "mult=hi")
		assert.Len(t, f, 2)
	})

	testutil.SwitchGlobal(&logger.OsExit, func(code int) {
		panic("osexit")
	})
	defer testutil.RestoreGlobals()

	t.Run("ErrColName", func(t *testing.T) {
		assert.PanicsWithValue(t, "osexit", func() {
			buildFormatterSetOrFail(accum, "bla=hi")
		})
	})

	t.Run("ErrBadExprName", func(t *testing.T) {
		assert.PanicsWithValue(t, "osexit", func() {
			buildFormatterSetOrFail(accum, "sum={unclosed")
		})
	})

	t.Run("ErrBadGlobalExpr", func(t *testing.T) {
		assert.PanicsWithValue(t, "osexit", func() {
			buildFormatterSetOrFail(accum, "{unclosed")
		})
	})

}
