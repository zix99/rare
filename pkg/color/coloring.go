package color

const (
	escapeCode     = "\x1b"
	foregroundCode = "[3"
)

type ColorCode string

const (
	Reset   ColorCode = escapeCode + "[0m"
	Red               = escapeCode + "[31m"
	Green             = escapeCode + "[32m"
	Yellow            = escapeCode + "[33m"
	Blue              = escapeCode + "[34m"
	Magenta           = escapeCode + "[35m"
	Cyan              = escapeCode + "[36m"
)

var groupColors = [...]string{Red, Green, Yellow, Blue, Magenta, Cyan}
