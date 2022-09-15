package cmd

import (
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
