package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	_ "honnef.co/go/tools/staticcheck"
)

func TestMain(t *testing.T) {
	assert.Error(t, cliMain("main"))
	assert.NoError(t, cliMain("main", "--help"))
}
