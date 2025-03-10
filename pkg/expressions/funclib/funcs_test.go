package funclib

import (
	"errors"
	"rare/pkg/expressions"
	"testing"

	"github.com/stretchr/testify/assert"
)

func voidFunc(args []expressions.KeyBuilderStage) (expressions.KeyBuilderStage, error) {
	return nil, nil
}

func TestFunctionSet(t *testing.T) {
	assert.NotZero(t, Builtins)
	assert.True(t, FunctionExists("sumi"))
	assert.False(t, FunctionExists("not-a-func"))
}

func TestAddFunction(t *testing.T) {
	AddFunctions(FunctionSet{
		"_test": voidFunc,
	})
	assert.True(t, FunctionExists("_test"))

	TryAddFunctions(FunctionSet{
		"_test": voidFunc,
	}, nil)
	TryAddFunctions(FunctionSet{
		"_test": voidFunc,
	}, errors.New("nope"))
}
