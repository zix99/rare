package expressions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextArray(t *testing.T) {
	ctx := KeyBuilderContextArray{
		Elements: []string{"a", "b"},
		Keys: map[string]string{
			"a": "b",
		},
	}
	assert.Equal(t, "a", ctx.GetMatch(0))
	assert.Equal(t, "b", ctx.GetMatch(1))
	assert.Equal(t, "", ctx.GetMatch(2))
	assert.Equal(t, "", ctx.GetKey("bla"))
	assert.Equal(t, "b", ctx.GetKey("a"))
}
