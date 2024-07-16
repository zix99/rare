package cmd

import (
	"os"
	"rare/pkg/expressions/funclib"
	"rare/pkg/testutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpressionCmd(t *testing.T) {
	testCommandSet(t, expressionCommand(),
		`--help`,
		`"a b c"`,
		`-b -s -d test --key a=b "abc {0} {a}"`,
	)
}

func TestExpressionOnlyOutput(t *testing.T) {
	o, e, err := testCommandCapture(expressionCommand(), `-d bob "hello {0}"`)
	assert.NoError(t, err)
	assert.Equal(t, e, "")
	assert.Equal(t, o, "hello bob\n")
}

func TestExpressionReadStdin(t *testing.T) {
	o, e, err := testutil.Capture(func(w *os.File) error {
		go func() {
			w.WriteString("hello {0}")
			w.Close()
		}()
		return testCommand(expressionCommand(), `-n -d bob -`)
	})
	assert.NoError(t, err)
	assert.Equal(t, e, "")
	assert.Equal(t, o, "hello bob")
}

func TestExpressionResults(t *testing.T) {
	o, e, err := testCommandCapture(expressionCommand(), `-s -d bob "abc {0}"`)
	assert.NoError(t, err)
	assert.Empty(t, e)

	assert.Equal(t,
		`Expression: abc {0}
Result:     abc bob

Stats
  Stages:        2
  Match Lookups: 1
  Key   Lookups: 0
`, o)
}

func TestExpressionErrors(t *testing.T) {
	o, e, err := testCommandCapture(expressionCommand(), "")
	assert.Error(t, err)
	assert.Empty(t, o)
	assert.NotEmpty(t, e)

	o, e, err = testCommandCapture(expressionCommand(), `-s ""`)
	assert.Error(t, err)
	assert.Empty(t, o)
	assert.NotEmpty(t, e)

	o, e, err = testCommandCapture(expressionCommand(), `-s "unterm {"`)
	assert.Error(t, err)
	assert.Empty(t, o)
	assert.NotEmpty(t, e)
}

func TestListFuncs(t *testing.T) {
	testutil.StoreGlobal(&funclib.Additional)
	defer testutil.RestoreGlobals()

	funclib.Additional["test"] = nil

	o, e, err := testCommandCapture(expressionCommand(), "--listfuncs")
	assert.NoError(t, err)
	assert.Empty(t, e)
	assert.Contains(t, o, "Builtin:")
	assert.Contains(t, o, "FuncsFile: test")
}

func TestKeyParser(t *testing.T) {
	k, v := parseKeyValue("")
	assert.Empty(t, k)
	assert.Empty(t, v)

	k, v = parseKeyValue("a")
	assert.Equal(t, "a", k)
	assert.Equal(t, "a", v)

	k, v = parseKeyValue("a=b")
	assert.Equal(t, "a", k)
	assert.Equal(t, "b", v)

	k, v = parseKeyValue("a=b=c")
	assert.Equal(t, "a", k)
	assert.Equal(t, "b=c", v)
}
