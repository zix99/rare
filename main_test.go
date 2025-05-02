package main

import (
	"rare/pkg/logger"
	"rare/pkg/testutil"
	"rare/pkg/testutil/acceptance"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	_ "honnef.co/go/tools/staticcheck"
)

func TestMain(t *testing.T) {
	assert.Error(t, cliMain("main"))
	assert.NoError(t, cliMain("main", "--help"))
}

func TestDupeCommands(t *testing.T) {
	commands := make(map[string]struct{})
	app := buildApp()

	for _, cmd := range app.Commands {
		for _, name := range cmd.Names() {
			assert.NotContains(t, commands, name)
			commands[name] = struct{}{}
		}
	}
}

func TestDupeFlags(t *testing.T) {
	app := buildApp()
	for _, cmd := range app.Commands {
		flags := make(map[string]struct{})

		// Global (even though technically can dupe, it's confusing, so prevent)
		for _, flag := range app.Flags {
			for _, name := range flag.Names() {
				assert.NotContains(t, flags, name)
				flags[name] = struct{}{}
			}
		}

		// And the specific flags
		for _, flag := range cmd.Flags {
			for _, name := range flag.Names() {
				assert.NotContains(t, flags, name)
				flags[name] = struct{}{}
			}
		}
	}
}

// Run all the tests in acceptance.tests
// See that file for format details
// Run only these tests with this command:
// go test -timeout 30s -run ^TestRunAcceptance$ ./ -v
func TestRunAcceptance(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	testutil.SwitchGlobal(&logger.OsExit, func(code int) {
		panic("exited code " + strconv.Itoa(code))
	})
	defer testutil.RestoreGlobals()

	acceptance.RunTestSuiteFile(t, "acceptance.tests", cliMain)
}
