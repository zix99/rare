package testutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatcher(t *testing.T) {
	assert.NoError(t, matchesPattern("hi", "hi"))
	assert.NoError(t, matchesPattern("hello *!", "hello bob!"))
	assert.NoError(t, matchesPattern("hello m?n", "hello mon"))
	assert.NoError(t, matchesPattern("hello m?n(", "hello mon("))

	assert.Error(t, matchesPattern("hi", "bye"))
	assert.Error(t, matchesPattern("hello m?n", "hello moon"))
	assert.Error(t, matchesPattern("hello *", "bye bob"))
	assert.Error(t, matchesPattern("hello *!", "hello bob?"))
}
