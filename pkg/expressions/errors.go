package expressions

const (
	ErrorBucket     = "<BUCKET-ERROR>" // Unable to bucket from given value (wrong type)
	ErrorBucketSize = "<BUCKET-SIZE>"  // Unable to get the size of the bucket (wrong type)
	ErrorType       = "<BAD-TYPE>"     // Error parsing the principle value of the input because of unexpected type
	ErrorParsing    = "<PARSE-ERROR>"  // Error parsing the principle value of the input
	ErrorArgCount   = "<ARGN>"         // Function to not support a variation with the given argument count
	ErrorArgName    = "<NAME>"         // A variable accessed by a given name does not exist
)
