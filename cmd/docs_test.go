package cmd

import "testing"

func TestDocs(t *testing.T) {
	testCommandSet(t, docsCommand(),
		``,
		`expressions`,
		`no-exist`,
	)
}
