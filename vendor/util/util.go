package util

import "os"

// Params support named-paaram
type Params = map[string](interface{})

// Identifier as a unique identifier, like ID
type Identifier string

var emptyIdentifier = *new(Identifier)

func (n Identifier) Empty() bool {
	return n == emptyIdentifier
}

// WorkingDir for agenda.
func WorkingDir() string {
	location, existed := os.LookupEnv("HOME")
	if !existed {
		location = "."
	}
	// NOTE: here to ensure workingdir existed ?
	ret := location + "/.agenda.d/"
	return ret
}
