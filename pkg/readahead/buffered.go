package readahead

import (
	"bytes"
	"io"
)

type BufferedReadAhead struct {
	r         io.Reader
	maxBufLen int

	buf    []byte
	offset int
	eof    bool

	token []byte
	delim byte

	onError func(error) // OnError is called if there are any downstream errors
}

var _ Scanner = &BufferedReadAhead{}

/*
Buffered read-ahead similar to Scanner, except it will leave the large-buffers in place
(rather than shifting them) so that a given slice is good for the duration of its life

This allows a slice reference to be passed around without worrying that the underlying data will change
which limits the amount the data needs to be copied around

Initial benchmarks shows a 8% savings over Scanner
*/
func NewBuffered(reader io.Reader, maxBufLen int) *BufferedReadAhead {
	if maxBufLen <= 1 {
		panic("Buf length must be > 1")
	}
	return &BufferedReadAhead{
		r:         reader,
		maxBufLen: maxBufLen,
		delim:     '\n',
	}
}

// Scan for the next token with a new line
func (s *BufferedReadAhead) Scan() bool {
	for {
		//var a chars
		relIndex := bytes.IndexByte(s.buf[s.offset:], s.delim)

		if relIndex >= 0 {
			start := s.offset
			s.offset += relIndex + 1
			s.token = dropCR(s.buf[start : start+relIndex])
			return true
		}

		// No new line, so either:
		// A) There's not enough room in the buffer
		// B) We're towards the end of buf, and need more data

		// Eof, so the rest of it will count as a line
		if s.eof && s.offset < len(s.buf) {
			ret := s.buf[s.offset:]
			s.offset = len(s.buf)
			s.token = ret
			return true
		} else if !s.eof {
			// Not enough in buffer to find next new-line.. need to fill until finding

			// Resize buffer, and copy over old buffer data
			oldbuf := s.buf
			s.buf = make([]byte, maxi(s.maxBufLen, len(oldbuf)-s.offset+s.maxBufLen/2))
			copy(s.buf, oldbuf[s.offset:])
			readOffset := len(oldbuf) - s.offset

			// Fill remaining buffer
			for readOffset < len(s.buf) {
				n, err := s.r.Read(s.buf[readOffset:])
				readOffset += n
				if err != nil {
					if err != io.EOF && s.onError != nil {
						s.onError(err)
					}
					s.eof = true
					break
				}
			}

			// Trim buffer to read-length, and reset offset
			s.buf = s.buf[:readOffset]
			s.offset = 0
		} else {
			s.token = nil
			return false
		}

	}
}

// Bytes retrieves the current bytes of the current token (line)
func (s *BufferedReadAhead) Bytes() []byte {
	return s.token
}

// ReadLine is shorthand for Scan() Token()
func (s *BufferedReadAhead) ReadLine() []byte {
	if !s.Scan() {
		return nil
	}
	return s.token
}

func (s *BufferedReadAhead) OnError(onError OnScannerError) {
	s.onError = onError
}
