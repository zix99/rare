package followreader

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/fsnotify.v1"
)

type NotifyFollowReader struct {
	filename string
	f        io.ReadSeekCloser

	OnError OnTailError

	watcherExit chan struct{}
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
		watcherExit: make(chan struct{}),
		eventWrite:  make(chan struct{}),
		eventDelete: make(chan struct{}, 1),
	}

	err = ret.startWatchEvents()
	if err != nil {
		f.Close()
		return nil, fmt.Errorf("unable to start notify: %w", err)
	}

	return ret, nil
}

func (s *NotifyFollowReader) Close() error {
	if s.f != nil {
		s.f.Close()
		s.f = nil
	}

	s.watcherExit <- struct{}{}

	return nil
}

func (s *NotifyFollowReader) Drain() error {
	_, err := s.f.Seek(0, os.SEEK_END)
	return err
}

func (s *NotifyFollowReader) Read(buf []byte) (int, error) {
	for {
		n, err := s.f.Read(buf)
		if err != nil && err != io.EOF {
			s.callOnError(err)
			return n, err
		}

		if n > 0 {
			return n, err
		}

		select {
		case <-s.eventDelete:
			return 0, io.EOF // TODO: Implement re-open
		case <-s.eventWrite:
		}
	}
}

func (s *NotifyFollowReader) startWatchEvents() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	if watcher.Add(s.filename) != nil {
		return err
	}

	s.watcherExit = make(chan struct{})
	s.eventWrite = make(chan struct{})
	s.eventDelete = make(chan struct{}, 1)

	go func() {
		defer watcher.Close()
		for {
			select {
			case <-s.watcherExit:
				return
			case event, ok := <-watcher.Events:
				if !ok {
					return
				} else if event.Op&fsnotify.Write != 0 {
					writeSignalNonBlock(s.eventWrite)
				} else if event.Op&fsnotify.Remove != 0 {
					writeSignalNonBlock(s.eventDelete)
				}
			}
		}
	}()

	return nil
}

func (s *NotifyFollowReader) callOnError(err error) {
	if s.OnError != nil {
		s.OnError(err)
	}
}

func writeSignalNonBlock(c chan<- struct{}) {
	select {
	case c <- struct{}{}:
	default:
	}
}
