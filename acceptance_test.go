package main

import (
	"bufio"
	"io"
	"os"
	"rare/pkg/testutil"
	"strconv"
	"strings"
	"testing"
)

type testConfig struct {
	name        string
	cmd         string
	stdout      strings.Builder
	stderr      strings.Builder
	outComp     stringComparer
	errComp     stringComparer
	expectError string
}

type stringComparer func(string, string) bool

var stringMatchers = map[string]stringComparer{
	"default":    strings.HasPrefix,
	"prefix":     strings.HasPrefix,
	"suffix":     strings.HasSuffix,
	"contains":   strings.Contains,
	"exact":      func(s1, s2 string) bool { return s1 == s2 },
	"ignorecase": strings.EqualFold,
}

// Run all the tests in acceptance.tests
// See that file for format details
// Run only these tests with this command:
// go test -timeout 30s -run ^TestRunAcceptance$ ./ -v
func TestRunAcceptance(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	f, err := os.Open("acceptance.tests")
	if err != nil {
		t.Fatal("Unable to read tests file")
	}
	defer f.Close()

	for test := range iterateTestDefinitions(t, f) {
		t.Run(test.name, func(t *testing.T) {
			runTestConfig(t, test)
		})
	}
}

func runTestConfig(t *testing.T, cfg testConfig) {
	t.Logf("RUN: %s", cfg.cmd)

	args := append([]string{"rare"}, testutil.SplitQuotedString(cfg.cmd)...)
	sout, serr, err := testutil.Capture(func(w *os.File) error {
		return cliMain(args...)
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

func iterateTestDefinitions(t *testing.T, r io.Reader) func(func(yield testConfig) bool) {
	scanner := bufio.NewScanner(r)

	return func(yield func(testConfig) bool) {
		var cfg testConfig
		var writeTarget *strings.Builder

		trimString := ""
		matcher := stringMatchers["prefix"]

	SCANNER:
		for scanner.Scan() {
			line := scanner.Text()

			switch {
			// Global switches
			case strings.HasPrefix(line, "#"): // skip
			case strings.HasPrefix(line, "INDENT "): // possible indent len for stdout/stderr
				indent, _ := strconv.Atoi(line[7:])
				trimString = strings.Repeat(" ", indent)
			case strings.HasPrefix(line, "MATCH "):
				matcher = stringMatchers[strings.ToLower(line[6:])]
				if matcher == nil {
					t.Fatalf("Unknown matcher: %s", line)
				}

			// Test definition
			case strings.HasPrefix(line, "NAME "):
				cfg.name = line[5:]
			case strings.HasPrefix(line, "RUN "):
				cfg.cmd = line[4:]
				cfg.outComp = matcher
				cfg.errComp = matcher
				writeTarget = &cfg.stdout
			case strings.HasPrefix(line, "STDOUT"):
				writeTarget = &cfg.stdout
				cfg.outComp = matcher
			case strings.HasPrefix(line, "STDERR"):
				writeTarget = &cfg.stderr
				cfg.errComp = matcher
			case strings.HasPrefix(line, "ERR "):
				cfg.expectError = line[4:]
			case writeTarget != nil && strings.HasPrefix(line, "END"): // execute
				writeTarget = nil

				if !yield(cfg) {
					break SCANNER
				}

				cfg = testConfig{} // reset
			case writeTarget != nil:
				writeTarget.WriteString(strings.TrimPrefix(line, trimString))
				writeTarget.WriteRune('\n')

			// Unknown line (non-blank)
			case line != "":
				t.Errorf("Unexpected line: %s", line)
			}
		}

		if scanner.Err() != nil {
			t.Errorf("Error parsing test file: %v", scanner.Err())
		}
	}
}
