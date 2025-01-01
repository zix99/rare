package dissect

import "errors"

var (
	ErrorKeyConflict     = errors.New("key conflict")
	ErrorUnclosedToken   = errors.New("unclosed token")
	ErrorSequentialToken = errors.New("sequential token")
)
