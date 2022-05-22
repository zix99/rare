package readahead

import (
	"bytes"
	"io"
)

type ImmediateReadAhead struct {
	r io.Reader

	buf     []byte
	bufSize int
	offset  int
	end     int

	delim byte

	token []byte
	eof   bool

	onError OnScannerError
}

var _ Scanner = &ImmediateReadAhead{}

/*
ReadThru is slightly different than Scanner and read-ahead given that:
 - Unlike Scanner and like ReadAhead, ReadThru leaves the buffer in-place
 - Like Scanner, and unlike ReadAhead, Scan() returns immediately after finding a delim,
   rather than blocking and buffering up to bufSize
*/
func NewImmediate(reader io.Reader, bufSize int) *ImmediateReadAhead {
	return &ImmediateReadAhead{
		r:       reader,
		bufSize: bufSize,
		buf:     make([]byte, bufSize),
		delim:   '\n',
	}
}

func (s *ImmediateReadAhead) Scan() bool {
RESTART:

	if s.offset < s.end {
		if eol := bytes.IndexByte(s.buf[s.offset:s.end], s.delim); eol >= 0 {
			s.token = dropCR(s.buf[s.offset : s.offset+eol])
			s.offset += eol + 1
			return true
		}
		if s.eof {
			s.token = s.buf[s.offset:s.end]
			s.offset = s.end
			return true
		}
	} else if s.eof {
		return false
	}

	// Read loop
	for {
		// Increase buf if needed (heuristically)
		if s.end >= len(s.buf) {
			old := s.buf
			s.buf = make([]byte, s.end-s.offset+s.bufSize)
			copy(s.buf, old[s.offset:s.end])
			s.end -= s.offset
			s.offset = 0
		}

		// Read data and check for errors
		n, err := s.r.Read(s.buf[s.end:])
		s.end += n

		if err != nil {
			s.eof = true
			if err != io.EOF && s.onError != nil {
				s.onError(err)
			}
			goto RESTART
		}

		// Check only the most recently read bytes for a new line
		if eol := bytes.IndexByte(s.buf[s.end-n:s.end], s.delim); eol >= 0 {
			end := s.end - n + eol
			s.token = dropCR(s.buf[s.offset:end])
			s.offset = end + 1
			return true
		}
	}
}

func (s *ImmediateReadAhead) Bytes() []byte {
	return s.token
}

func (s *ImmediateReadAhead) ReadLine() []byte {
	if s.Scan() {
		return s.token
	}
	return nil
}

func (s *ImmediateReadAhead) OnError(f OnScannerError) {
	s.onError = f
}
