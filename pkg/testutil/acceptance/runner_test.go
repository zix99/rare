package acceptance

import (
	"fmt"
	"os"
	"testing"
)

func TestSimpleTests(t *testing.T) {
	RunTestSuiteFile(t, "example.tests", func(args ...string) error {
		fmt.Printf("Expected '%s' and '%s'\n", args[1], args[2])
		fmt.Fprintln(os.Stderr, "err")
		return nil
	})
}
