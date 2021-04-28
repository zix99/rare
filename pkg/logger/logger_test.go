package logger

import "testing"

func TestLog(t *testing.T) {
	Println("Howdy")
	DeferLogs()
	Println("Buffered")
	ImmediateLogs()
}
