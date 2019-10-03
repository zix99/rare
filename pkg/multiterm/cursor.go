package multiterm

import "fmt"

func escape(format string, args ...interface{}) string {
	const ESCAPE = "\x1b"
	return ESCAPE + fmt.Sprintf(format, args...)
}

func moveCursorf(line, col int) string {
	return escape("[%d;%dH", line, col)
}

func moveCursor(line, col int) {
	fmt.Print(moveCursorf(line, col))
}

func moveUpf(n int) string {
	return escape("[%dA", n)
}

func moveUp(n int) {
	fmt.Print(moveUpf(n))
}
