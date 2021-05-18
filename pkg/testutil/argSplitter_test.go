package testutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArgSplit(t *testing.T) {
	args := SplitQuotedString(`this is "a test" "" "thing" "noquote other`)
	assert.Equal(t, []string{"this", "is", "a test", "", "thing", "noquote", "other"}, args)
}
