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

func TestFilterExtract(t *testing.T) {
	out, eout, err := testCommandCapture(filterCommand(), `-m (\d+) -e "{1}" testdata/log.txt`)
	assert.NoError(t, err)
	assert.Equal(t, "5\n22\n5\n", out)
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
