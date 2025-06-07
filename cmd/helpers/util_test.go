package helpers

import (
	"testing"

	"github.com/zix99/rare/pkg/logger"

	"github.com/stretchr/testify/assert"
)

func testLogFatal(t *testing.T, expectsCode int, f func()) (code int) {
	code = -1

	oldExit := logger.OsExit
	defer func() {
		logger.OsExit = oldExit
	}()
	logger.OsExit = func(v int) {
		code = v
		panic("logger.osexit")
	}

	assert.PanicsWithValue(t, "logger.osexit", f)
	assert.Equal(t, expectsCode, code)
	return
}
