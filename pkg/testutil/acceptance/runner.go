package acceptance

import (
	"errors"
	"io"
	"os"
	"rare/pkg/testutil"
	"strings"
	"testing"
)

type Runner func(args ...string) error

func RunTestSuite(t *testing.T, r io.Reader, runner Runner) {
	for test := range iterateTestDefinitions(t, r) {
		t.Run(test.name, func(t *testing.T) {
			runTestConfig(t, &test, runner)
		})
	}
}

func RunTestSuiteFile(t *testing.T, filename string, runner Runner) {
	f, err := os.Open(filename)
	if err != nil {
		t.Fatalf("Unable to open test file: %s", filename)
	}
	defer f.Close()

	RunTestSuite(t, f, runner)
}

func runTestConfig(t *testing.T, cfg *testConfig, runner Runner) {
	t.Logf("RUN (line %d): %s", cfg.linenum, cfg.cmd)

	args := append([]string{"app"}, testutil.SplitQuotedString(cfg.cmd)...)
	sout, serr, err := testutil.Capture(func(w *os.File) (ret error) {
		if cfg.stdin.Len() > 0 {
			t.Logf("Copying %d bytes to stdin", cfg.stdin.Len())
			go func() {
				io.Copy(w, strings.NewReader(cfg.stdin.String()))
				w.Close()
			}()
		}

		// catch panics
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Recovered from panic: %v", r)
				ret = errors.New(r.(string))
			}
		}()

		return runner(args...)
	})

	if (err != nil && err.Error() != cfg.expectError) || (err == nil && cfg.expectError != "") {
		t.Errorf("ERROR (line %d): '%v', expected '%s'", cfg.linenum, err, cfg.expectError)
	}

	if !cfg.outComp(sout, cfg.stdout.String()) {
		t.Errorf("STDOUT (line %d) Expected:\n%s\nGot:\n%s\n", cfg.linenum, cfg.stdout.String(), sout)
	}

	if !cfg.errComp(serr, cfg.stderr.String()) || (len(serr) > 0 && cfg.stderr.Len() == 0) {
		t.Errorf("STDERR (line %d) Expected:\n%s\nGot:\n%s\n", cfg.linenum, cfg.stderr.String(), serr)
	}

	t.Logf("DONE: err=%v; stderr=%s; len(stdout)=%d", err, serr, len(sout))
}
