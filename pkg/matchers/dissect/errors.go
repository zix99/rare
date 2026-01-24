package dissect

import "errors"

var (
	ErrorKeyConflict     = errors.New("dissect error: key conflict")
	ErrorUnclosedToken   = errors.New("dissect error: unclosed token")
	ErrorSequentialToken = errors.New("dissect error: sequential token")
)
