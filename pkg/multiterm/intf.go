package multiterm

type MultilineTerm interface {
	WriteForLine(line int, s string)
	Close()
}
