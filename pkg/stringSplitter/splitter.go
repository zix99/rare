package stringSplitter

import "strings"

/*
A splitter without making heap memory
*/

type Splitter struct {
	S     string
	Delim string
	next  int
}

func (s *Splitter) Next() (ret string) {
	if s.next < 0 {
		return ""
	}

	idx := strings.Index(s.S[s.next:], s.Delim)
	if idx < 0 {
		ret = s.S[s.next:]
		s.next = -1
		return
	}
	idx += s.next

	ret = s.S[s.next:idx]
	s.next = idx + len(s.Delim)
	return
}

func (s *Splitter) NextOk() (ret string, ok bool) {
	ok = !s.Done()
	ret = s.Next()
	return
}

func (s *Splitter) Done() bool {
	return s.next < 0
}
