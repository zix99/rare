package multiterm

import "testing"

func TestBasicHisto(t *testing.T) {
	mt := NewHistogram(5)
	mt.WriteForLine(4, "key", 1000)
}
