package termrenderers

import (
	"testing"

	"github.com/zix99/rare/pkg/multiterm"
	"github.com/zix99/rare/pkg/multiterm/termunicode"

	"github.com/stretchr/testify/assert"
)

func TestBasicHisto(t *testing.T) {
	vt := multiterm.NewVirtualTerm()
	mt := NewHistogram(vt, 5)

	termunicode.UnicodeEnabled = false
	mt.UpdateTotal(10000)
	mt.WriteForLine(1, "key", 1000)
	mt.WriteFooter(0, "hello")
	termunicode.UnicodeEnabled = true

	assert.Equal(t, "", vt.Get(0))
	assert.Equal(t, "key         1,000      [10.0%] ||||||||||||||||||||||||||||||||||||||||||||||||||", vt.Get(1))
	assert.Equal(t, "", vt.Get(2))
	assert.Equal(t, "", vt.Get(3))
	assert.Equal(t, "", vt.Get(4))
	assert.Equal(t, "hello", vt.Get(5))

	mt.Close()
	assert.True(t, vt.IsClosed())
}

func TestHistoWithValues(t *testing.T) {
	vt := multiterm.NewVirtualTerm()
	mt := NewHistogram(vt, 2)
	mt.ShowBar = false
	mt.WriteForLine(0, "abc", 1234)
	mt.WriteForLine(1, "q", 123)
	mt.WriteForLine(3, "abn", 444)
	assert.Equal(t, "abc         1,234     ", vt.Get(0))
	assert.Equal(t, "q           123       ", vt.Get(1))
	assert.Equal(t, "", vt.Get(2))
}

func TestHistoWithUnicode(t *testing.T) {
	vt := multiterm.NewVirtualTerm()
	mt := NewHistogram(vt, 2)
	mt.ShowBar = false
	mt.textSpacing = 1

	mt.WriteForLine(0, "✤✥✦", 1)
	mt.WriteForLine(1, "abc", 1)

	assert.Equal(t, "✤✥✦    1         ", vt.Get(0))
	assert.Equal(t, "abc    1         ", vt.Get(1))
	assert.Equal(t, "", vt.Get(2))
}
