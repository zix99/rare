package stdlib

import (
	"errors"
	"fmt"

	. "github.com/zix99/rare/pkg/expressions" //lint:ignore ST1001 Legacy
)

type funcError struct {
	expr string
	err  error
}

func newFuncErr(expr, message string) funcError {
	return funcError{expr, errors.New(message)}
}

// Realtime errors
const (
	ErrorNum      = "<BAD-TYPE>"    // Error parsing the principle value of the input because of unexpected type (numeric)
	ErrorParsing  = "<PARSE-ERROR>" // Error parsing the principle value of the input (non-numeric)
	ErrorArgCount = "<ARGN>"        // Function to not support a variation with the given argument count
	ErrorConst    = "<CONST>"       // Expected constant value
	ErrorEnum     = "<ENUM>"        // A given value is not contained within a set
	ErrorArgName  = "<NAME>"        // A variable accessed by a given name does not exist
	ErrorEmpty    = "<EMPTY>"       // A value was expected, but was empty
	ErrorFile     = "<FILE>"        // Unable to read file
	ErrorValue    = "<VALUE>"       // Value is out of range or invalid (eg. range incrementer is 0)
)

// Compilation errors
var (
	ErrNum     = newFuncErr(ErrorNum, "invalid arg type, expected int") // always numeric
	ErrParsing = newFuncErr(ErrorParsing, "unable to parse")            // always non-numeric
	ErrConst   = newFuncErr(ErrorConst, "expected const")
	ErrEnum    = newFuncErr(ErrorEnum, "unable to find value in set")
	ErrEmpty   = newFuncErr(ErrorEmpty, "invalid empty value")
	ErrFile    = newFuncErr(ErrorFile, "unable to read file")
	ErrValue   = newFuncErr(ErrorValue, "value out of range")
)

var (
	ErrArgCount = errors.New("invalid number of arguments")
)

func stageError(err funcError) (KeyBuilderStage, error) {
	return func(ctx KeyBuilderContext) string {
		return err.expr
	}, err.err
}

func stageErrorf(err funcError, msg string) (KeyBuilderStage, error) {
	return func(ctx KeyBuilderContext) string {
		return err.expr
	}, fmt.Errorf("%s, %w", msg, err.err)
}

func stageArgError(err funcError, argIndex int) (KeyBuilderStage, error) {
	return func(ctx KeyBuilderContext) string {
		return err.expr
	}, fmt.Errorf("argument %d, %w", argIndex+1, err.err)
}

func stageErrArgCount(got []KeyBuilderStage, expected int) (KeyBuilderStage, error) {
	return stageError(funcError{
		ErrorArgCount,
		fmt.Errorf("%w: got %d, expected %d", ErrArgCount, len(got), expected),
	})
}

func stageErrArgRange(got []KeyBuilderStage, text string) (KeyBuilderStage, error) {
	return stageError(funcError{
		ErrorArgCount,
		fmt.Errorf("%w: got %d, expected %s", ErrArgCount, len(got), text),
	})
}
