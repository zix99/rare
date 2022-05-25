package followreader

import (
	"fmt"
	"io"
	"os"
	"time"
)

type PollingTailReader struct {
	filename string
	f        io.ReadSeekCloser

	// Metrics
	readBytes int64

	// Options
	ReadAttempts int           // Number of read-attempts before checking to re-open
	PollDelay    time.Duration // Delay between read attempts
	Reopen       bool          // If true, will try to open() file again if it looks different after a delay; false will send EOF if file goes away
	OnError      OnTailError   // Called on unhandled error
}

var _ FollowReader = &PollingTailReader{}

func NewPolling(filename string, reopen bool) (*PollingTailReader, error) {
	f, err := os.Open(filename)

	if err != nil && !reopen {
		return nil, fmt.Errorf("unable to open file and cannot reopen: %w", err)
	}

	ret := &PollingTailReader{
		filename:     filename,
		f:            f,
		PollDelay:    250 * time.Millisecond,
		ReadAttempts: 5,
		Reopen:       reopen,
	}

	return ret, nil
}

// Drain navigates to the end of the stream
func (s *PollingTailReader) Drain() error {
	_, err := s.f.Seek(0, os.SEEK_END)
	return err
}

// Close file and underlying resources
func (s *PollingTailReader) Close() error {
	if s.f != nil {
		s.f.Close()
		s.f = nil
	}

	return nil
}

func (s *PollingTailReader) Read(buf []byte) (int, error) {
	for {
		if s.f != nil {
			for i := 0; i < s.ReadAttempts; i++ {
				n, err := s.f.Read(buf)
				s.readBytes += int64(n)

				if n > 0 {
					return n, nil
				}

				if err != nil && err != io.EOF {
					// Any error other than EOF is raised
					s.callOnError(err)
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
					// Assume new file
					if seeker, ok := s.f.(io.Seeker); ok {
						seeker.Seek(s.readBytes, io.SeekStart)
					}
				} else {
					s.readBytes = 0
				}
			}
		} else { // No re-open, if the file's missing, that's EOF
			_, err := os.Stat(s.filename)
			if err != nil {
				return 0, io.EOF
			}
		}

	}
}

func (s *PollingTailReader) callOnError(err error) {
	if s.OnError != nil {
		s.OnError(err)
	}
}
