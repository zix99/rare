package followreader

import (
	"fmt"
	"io"
	"os"
	"time"
)

type PollingFollowReader struct {
	filename string
	f        *os.File

	// State
	readBytes int64
	closed    bool

	// Options
	ReadAttempts int           // Number of read-attempts before checking to re-open
	PollDelay    time.Duration // Delay between read attempts
	Reopen       bool          // If true, will try to open() file again if it looks different after a delay; false will send EOF if file goes away
}

var _ FollowReader = &PollingFollowReader{}

func NewPolling(filename string, reopen bool) (*PollingFollowReader, error) {
	f, err := os.Open(filename)

	if err != nil && !reopen {
		return nil, fmt.Errorf("unable to open file and cannot reopen: %w", err)
	}

	ret := &PollingFollowReader{
		filename:     filename,
		f:            f,
		PollDelay:    250 * time.Millisecond,
		ReadAttempts: 5,
		Reopen:       reopen,
	}

	return ret, nil
}

// Drain navigates to the end of the stream
func (s *PollingFollowReader) Drain() error {
	offset, err := s.f.Seek(0, os.SEEK_END)
	if err == nil {
		s.readBytes = offset
	}
	return err
}

// Close file and underlying resources
func (s *PollingFollowReader) Close() error {
	if s.f != nil {
		s.f.Close()
		s.f = nil
	}

	s.closed = true

	return nil
}

func (s *PollingFollowReader) Read(buf []byte) (int, error) {
	if s.closed {
		return 0, io.EOF
	}

	for {
		if s.f != nil {
			for i := 0; i < s.ReadAttempts; i++ {
				n, err := s.f.Read(buf)
				s.readBytes += int64(n)

				if n > 0 {
					return n, nil
				}

				if err != nil && err != io.EOF {
					return n, err
				}

				time.Sleep(s.PollDelay)
			}
		} else {
			time.Sleep(s.PollDelay)
		}

		// Didn't read any bytes... has the file inode changed?
		if s.Reopen {
			st, _ := os.Stat(s.filename)
			if st != nil && st.Size() != s.readBytes {
				s.f, _ = os.Open(s.filename)
				if st.Size() >= s.readBytes {
					// Likely existing file that is re-opened, start reading where we left off
					s.f.Seek(s.readBytes, io.SeekStart)
				} else {
					// Assume new file, restart reading from beginning
					s.readBytes = 0
				}
			}
		} else { // No re-open, if the file's missing, that's EOF
			_, err := os.Stat(s.filename)
			if err != nil {
				s.Close()
				return 0, io.EOF
			}
		}

	}
}
