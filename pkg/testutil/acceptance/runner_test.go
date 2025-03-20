package acceptance

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

func TestSimpleTests(t *testing.T) {
	RunTestSuiteFile(t, "example.tests", func(args ...string) error {
		if args[1] == "error" {
			return errors.New("failed")
		}
		fmt.Printf("Expected '%s' and '%s'\n", args[1], args[2])
		fmt.Fprintln(os.Stderr, "err")
		return nil
	})
}

func TestStdinTests(t *testing.T) {
	RunTestSuiteFile(t, "stdin.tests", func(args ...string) error {
		var s string
		fmt.Scanln(&s)
		fmt.Printf("Expected '%s' and '%s' with stdin '%s'\n", args[1], args[2], s)
		return nil
	})
}
