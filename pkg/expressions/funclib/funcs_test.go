package funclib

import (
	"errors"
	"testing"

	"github.com/zix99/rare/pkg/expressions"

	"github.com/stretchr/testify/assert"
)

func voidFunc(args []expressions.KeyBuilderStage) (expressions.KeyBuilderStage, error) {
	return nil, nil
}

func TestFunctionSet(t *testing.T) {
	assert.NotZero(t, Builtins)
}

func TestAddFunction(t *testing.T) {
	AddFunctions(FunctionSet{
		"_test": voidFunc,
	})
	TryAddFunctions(FunctionSet{
		"_test": voidFunc,
	}, nil)
	TryAddFunctions(FunctionSet{
		"_test": voidFunc,
	}, errors.New("nope"))
}
