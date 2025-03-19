package acceptance

import (
	"io"
	"os"
	"rare/pkg/testutil"
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
	t.Logf("RUN: %s", cfg.cmd)

	args := append([]string{"rare"}, testutil.SplitQuotedString(cfg.cmd)...)
	sout, serr, err := testutil.Capture(func(w *os.File) error {
		return runner(args...)
	})

	if (err != nil && err.Error() != cfg.expectError) || (err == nil && cfg.expectError != "") {
		t.Errorf("ERROR: '%v', expected '%s'", err, cfg.expectError)
	}

	if !cfg.outComp(sout, cfg.stdout.String()) {
		t.Errorf("STDOUT Expected:\n%s\nGot:\n%s\n", cfg.stdout.String(), sout)
	}

	if !cfg.errComp(serr, cfg.stderr.String()) || (len(serr) > 0 && cfg.stderr.Len() == 0) {
		t.Errorf("STDERR Expected:\n%s\nGot:\n%s\n", cfg.stderr.String(), serr)
	}

	t.Logf("DONE: err=%v; stderr=%s; len(stdout)=%d", err, serr, len(sout))
}
