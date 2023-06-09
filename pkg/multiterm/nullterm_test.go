package multiterm

import "testing"

// null test for a null term (yay code coverage)
func TestNullTerm(t *testing.T) {
	nt := &NullTerm{}
	nt.WriteForLine(0, "bla")
	nt.WriteForLinef(1, "bla")
	nt.Close()
}
