package stdlib

const (
	ErrorType     = "<BAD-TYPE>"    // Error parsing the principle value of the input because of unexpected type
	ErrorParsing  = "<PARSE-ERROR>" // Error parsing the principle value of the input (non-numeric)
	ErrorArgCount = "<ARGN>"        // Function to not support a variation with the given argument count
	ErrorConst    = "<CONST>"       // Expected constant value
	ErrorEnum     = "<ENUM>"        // A given value is not contained within a set
	ErrorArgName  = "<NAME>"        // A variable accessed by a given name does not exist
	ErrorEmpty    = "<EMPTY>"       // A value was expected, but was empty
)
