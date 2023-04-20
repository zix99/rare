package stdlib

import (
	"errors"
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
	ErrorType    = "<BAD-TYPE>"    // Error parsing the principle value of the input because of unexpected type
	ErrorParsing = "<PARSE-ERROR>" // Error parsing the principle value of the input (non-numeric)
	// ErrorArgCount = "<ARGN>"        // Function to not support a variation with the given argument count
	// ErrorConst    = "<CONST>"       // Expected constant value
	// ErrorEnum     = "<ENUM>"        // A given value is not contained within a set
	ErrorArgName = "<NAME>" // A variable accessed by a given name does not exist
	// ErrorEmpty    = "<EMPTY>"       // A value was expected, but was empty
)

// Compilation errors
var (
	ErrTypeInt  = newFuncErr("<BAD-TYPE>", "invalid arg type, expected int") // always numeric?
	ErrParsing  = newFuncErr("<PARSE-ERROR>", "unable to parse")             // always non-numeric?
	ErrArgCount = newFuncErr("<ARGN>", "invalid number of arguments")
	ErrConst    = newFuncErr("<CONST>", "expected const")
	ErrEnum     = newFuncErr("<ENUM>", "unable to find value in set")
	ErrEmpty    = newFuncErr("<EMPTY>", "invalid empty value")
)
