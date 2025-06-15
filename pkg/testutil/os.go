package testutil

import (
	"runtime"
	"testing"
)

func IsWindows() bool {
	return runtime.GOOS == "windows"
}

func SkipWindows(t *testing.T) {
	if IsWindows() {
		t.Skip("skip windows")
	}
}
