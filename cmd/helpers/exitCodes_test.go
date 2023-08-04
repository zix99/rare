package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockExitState struct {
	batchErr, aggErr, extSum int
}

func (s *mockExitState) ReadErrors() int {
	return s.batchErr
}

func (s *mockExitState) ParseErrors() uint64 {
	return uint64(s.aggErr)
}

func (s *mockExitState) MatchedLines() uint64 {
	return uint64(s.extSum)
}

func TestDetermineErrorState(t *testing.T) {
	s := mockExitState{0, 0, 1}
	assert.NoError(t, DetermineErrorState(&s, &s, &s))

	s = mockExitState{0, 0, 0}
	assert.Error(t, DetermineErrorState(&s, &s, &s))

	s = mockExitState{0, 1, 1}
	assert.Error(t, DetermineErrorState(&s, &s, &s))

	s = mockExitState{1, 0, 1}
	assert.Error(t, DetermineErrorState(&s, &s, &s))
}
