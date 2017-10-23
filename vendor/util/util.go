package util

import (
	"log"
	"os"
)

var (
	Log  = log.Println
	Logf = log.Printf
	// Log = func(args ...interface{}) {}
	// Logf = func(args ...interface{}) {}
)

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
