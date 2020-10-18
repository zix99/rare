package fuzzy

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleTable(t *testing.T) {
	tbl := NewFuzzyTable(0.5, 5, 100)
	_, new := tbl.GetMatchId("test")
	assert.True(t, new)

	_, new = tbl.GetMatchId("test")
	assert.False(t, new)

	_, new = tbl.GetMatchId("blah")
	assert.True(t, new)

	_, new = tbl.GetMatchId("tast")
	assert.False(t, new)
}

func BenchmarkSimpleTable(b *testing.B) {
	tbl := NewFuzzyTable(0.7, 5, 100)
	for n := 0; n < b.N; n++ {
		tbl.GetMatchId(fmt.Sprintf("abcd-%d", n%100))
	}
}
