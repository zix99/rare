package cmd

import (
	"os"
	"rare/pkg/testutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestFilter(t *testing.T) {
	testCommandSet(t, filterCommand(),
		`-m \d+ testdata/log.txt`,
		`-m (\d+) testdata/log.txt`,
	)
}

func TestSarch(t *testing.T) {
	testCommandSet(t, searchCommand(),
		`the testdata/log.txt`,
		`testtest`,
	)
}

func TestSearchOutput(t *testing.T) {
	out, eout, err := testCommandCapture(searchCommand(), `last testdata/*.txt`)
	assert.NoError(t, err)
	assert.Equal(t, "testdata/log.txt 3: 5 is the last\n", out)
	assert.Equal(t, "Read   : 5 file(s) (9.41 KB)\nMatched: 1 / 76\n", eout)
}

func TestFilterExtract(t *testing.T) {
	out, eout, err := testCommandCapture(filterCommand(), `-m (\d+) -e "{1}" testdata/log.txt`)
	assert.NoError(t, err)
	assert.Equal(t, "5\n22\n5\n", out)
	assert.Equal(t, "Matched: 3 / 6\n", eout)
}

func TestFilterLine(t *testing.T) {
	out, eout, err := testCommandCapture(filterCommand(), `-l -m (\d+) -e "{1}" testdata/log.txt`)
	assert.NoError(t, err)
	assert.Equal(t, "testdata/log.txt 1: 5\ntestdata/log.txt 2: 22\ntestdata/log.txt 3: 5\n", out)
	assert.Equal(t, "Matched: 3 / 6\n", eout)
}

func TestFilterMultiExtract(t *testing.T) {
	out, eout, err := testCommandCapture(filterCommand(), `-m (\d+) -e "{1}" -e "b-{1}" testdata/log.txt`)
	assert.NoError(t, err)
	assert.Equal(t, "5\tb-5\n22\tb-22\n5\tb-5\n", out)
	assert.Equal(t, "Matched: 3 / 6\n", eout)
}

func TestFilterExtractFull(t *testing.T) {
	out, eout, err := testCommandCapture(filterCommand(), `-m (\d+) --extract "{1}" testdata/log.txt`)
	assert.NoError(t, err)
	assert.Equal(t, "5\n22\n5\n", out)
	assert.Equal(t, "Matched: 3 / 6\n", eout)
}

func TestFilterWithDissect(t *testing.T) {
	out, eout, err := testCommandCapture(filterCommand(), `-d "%{w0} is %{w1} " -e "{0}: {w0}={2}" testdata/log.txt`)
	assert.NoError(t, err)
	assert.Equal(t, "this is a : this=a\n22 is the : 22=the\n5 is the : 5=the\n", out)
	assert.Equal(t, "Matched: 3 / 6\n", eout)
}

func TestFilterWithDissectIgnoreCase(t *testing.T) {
	out, eout, err := testCommandCapture(filterCommand(), `-I -d "%{w0} IS %{w1} " -e "{0}: {w0}={2}" testdata/log.txt`)
	assert.NoError(t, err)
	assert.Equal(t, "this is a : this=a\n22 is the : 22=the\n5 is the : 5=the\n", out)
	assert.Equal(t, "Matched: 3 / 6\n", eout)
}

func TestFilterFromStdin(t *testing.T) {
	out, eout, err := testutil.Capture(func(w *os.File) error {
		go func() {
			w.WriteString("line 1\n")
			w.WriteString("line 5\n")
			w.WriteString("no number\n")
			w.Close()
		}()
		return testCommand(filterCommand(), `-m (\d+) -e "{src}:{line} {1}-{1}"`)
	})

	assert.NoError(t, err)
	assert.Equal(t, "<stdin>:1 1-1\n<stdin>:2 5-5\n", out)
	assert.Equal(t, "Matched: 2 / 3\n", eout)
}

func TestFilterFileNotExist(t *testing.T) {
	out, eout, err := testCommandCapture(filterCommand(), `-m (\d+) -e "{1}" testdata/no-exist.txt`)
	assert.Error(t, err)
	assert.Equal(t, 2, err.(cli.ExitCoder).ExitCode())
	assert.Equal(t, "", out)
	assert.Equal(t, "Matched: 0 / 0\nRead errors", eout)
}

func TestFilterNoMatches(t *testing.T) {
	out, eout, err := testCommandCapture(filterCommand(), `-m notfound(\d+) -e "{1}" testdata/log.txt`)
	assert.Error(t, err)
	assert.Equal(t, 1, err.(cli.ExitCoder).ExitCode())
	assert.Equal(t, "", out)
	assert.Equal(t, "Matched: 0 / 6\n", eout)
}
