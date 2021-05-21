package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	assert.Error(t, cliMain("main"))
	assert.NoError(t, cliMain("main", "--help"))
}
