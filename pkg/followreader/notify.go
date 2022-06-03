package followreader

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/fsnotify/fsnotify"
)

type NotifyFollowReader struct {
	filename string
	f        *os.File

	ReOpen bool

	closed      bool
	watcher     *fsnotify.Watcher
	eventWrite  chan struct{}
	eventDelete chan struct{}
}

var _ FollowReader = &NotifyFollowReader{}

func NewNotify(filename string, reopen bool) (*NotifyFollowReader, error) {
	f, err := os.Open(filename)

	if err != nil && !reopen {
		return nil, fmt.Errorf("unable to open file and cannot reopen: %w", err)
	}

	ret := &NotifyFollowReader{
		filename:    filename,
		f:           f,
		ReOpen:      reopen,
		eventWrite:  make(chan struct{}, 1),
		eventDelete: make(chan struct{}, 1),
	}

	ret.watcher, err = ret.startWatcher()
	if err != nil {
		if f != nil {
			f.Close()
		}
		return nil, fmt.Errorf("unable to start notify: %w", err)
	}

	return ret, nil
}

func (s *NotifyFollowReader) Close() error {
	if !s.closed {
		s.closeFile()
		s.watcher.Close()

		s.closed = true
	}

	return nil
}

func (s *NotifyFollowReader) Drain() error {
	if s.f != nil {
		_, err := s.f.Seek(0, os.SEEK_END)
		return err
	}
	return nil
}

func (s *NotifyFollowReader) Read(buf []byte) (int, error) {
	if s.closed {
		return 0, io.EOF
	}

	for {
		if s.f != nil {
			n, err := s.f.Read(buf)

			if n > 0 {
				return n, nil
			}
			if err != nil && err != io.EOF {
				return n, err
			}
		}

		// Wait for changes
		select {
		case <-s.eventWrite:
			if s.f == nil && s.ReOpen { // Re-open if able and willing
				if f, err := os.Open(s.filename); err == nil {
					s.f = f
				}
			}
		case <-s.eventDelete:
			if s.ReOpen {
				s.closeFile()
			} else {
				s.Close()
				return 0, io.EOF
			}
		}
	}
}

func (s *NotifyFollowReader) startWatcher() (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	if err := watcher.Add(path.Dir(s.filename)); err != nil {
		watcher.Close()
		return nil, err
	}

	go func() {
		defer watcher.Close()
		for {
			event, ok := <-watcher.Events
			switch {
			case !ok:
				return
			case path.Base(s.filename) != path.Base(event.Name):
				// nop
			case event.Op&fsnotify.Write != 0:
				writeSignalNonBlock(s.eventWrite)
			case event.Op&fsnotify.Remove != 0:
				writeSignalNonBlock(s.eventDelete)
			case event.Op&fsnotify.Create != 0:
				writeSignalNonBlock(s.eventWrite)
			}
		}
	}()

	return watcher, nil
}

func (s *NotifyFollowReader) closeFile() {
	if s.f != nil {
		s.f.Close()
		s.f = nil
	}
}

func writeSignalNonBlock(c chan<- struct{}) {
	select {
	case c <- struct{}{}:
	default:
	}
}
