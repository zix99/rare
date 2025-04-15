package acceptance

import (
	"bufio"
	"io"
	"strconv"
	"strings"
	"testing"
)

type stringComparer func(string, string) bool

var stringMatchers = map[string]stringComparer{
	"default":    strings.HasPrefix,
	"prefix":     strings.HasPrefix,
	"suffix":     strings.HasSuffix,
	"contains":   strings.Contains,
	"exact":      func(s1, s2 string) bool { return s1 == s2 },
	"ignorecase": strings.EqualFold,
}

type testConfig struct {
	name        string
	cmd         string
	stdout      strings.Builder
	stderr      strings.Builder
	stdin       strings.Builder
	outComp     stringComparer
	errComp     stringComparer
	expectError string
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
			case strings.HasPrefix(line, "STDIN"):
				writeTarget = &cfg.stdin
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
