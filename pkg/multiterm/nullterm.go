package multiterm

type NullTerm struct{}

var _ MultilineTerm = &NullTerm{}

func (s *NullTerm) WriteForLine(line int, l string) {}

func (s *NullTerm) WriteForLinef(line int, format string, args ...interface{}) {}

func (s *NullTerm) Close() {}
