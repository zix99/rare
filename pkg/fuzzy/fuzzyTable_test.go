package fuzzy

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleTable(t *testing.T) {
	tbl := NewFuzzyTable(0.5, 5)
	id, _, new := tbl.GetMatchId("test")
	assert.Equal(t, 0, id)
	assert.True(t, new)

	id, _, new = tbl.GetMatchId("test")
	assert.Equal(t, 0, id)
	assert.False(t, new)

	id, _, new = tbl.GetMatchId("blah")
	assert.Equal(t, 1, id)
	assert.True(t, new)

	id, _, new = tbl.GetMatchId("tast")
	assert.Equal(t, 0, id)
	assert.False(t, new)
}

func BenchmarkSimpleTable(b *testing.B) {
	tbl := NewFuzzyTable(0.7, 5)
	for n := 0; n < b.N; n++ {
		tbl.GetMatchId(fmt.Sprintf("abcd-%d", n%100))
	}
}
