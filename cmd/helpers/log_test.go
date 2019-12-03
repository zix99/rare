package helpers

import "testing"

func TestLog(t *testing.T) {
	ErrLog.Println("Howdy")
	EnableLogBuffer()
	ErrLog.Println("Buffered")
	DisableAndFlushLogBuffer()
}
