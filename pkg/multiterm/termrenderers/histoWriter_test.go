package termrenderers

import (
	"rare/pkg/multiterm"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicHisto(t *testing.T) {
	mt := NewHistogram(multiterm.New(), 5)
	mt.WriteForLine(4, "key", 1000)
}

func TestHistoWithValues(t *testing.T) {
	vt := multiterm.NewVirtualTerm()
	mt := NewHistogram(vt, 2)
	mt.ShowBar = false
	mt.WriteForLine(0, "abc", 1234)
	mt.WriteForLine(1, "q", 123)
	mt.WriteForLine(3, "abn", 444)
	assert.Equal(t, "abc                 1,234     ", vt.Get(0))
	assert.Equal(t, "q                   123       ", vt.Get(1))
	assert.Equal(t, "", vt.Get(2))
}
