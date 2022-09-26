package multiterm

type MultilineTerm interface {
	WriteForLine(line int, s string)
	WriteForLinef(line int, format string, args ...interface{})
	Close()
}
