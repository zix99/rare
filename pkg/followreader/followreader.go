package followreader

import (
	"io"
)

type FollowReader interface {
	io.ReadCloser
	Drain() error
}

func New(filename string, reopen, poll bool) (FollowReader, error) {
	if poll {
		return NewPolling(filename, reopen)
	}
	return NewNotify(filename, reopen)
}
