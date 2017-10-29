package util

import (
	"log"
	"os"
)

var (
	Log  = log.Println
	Logf = log.Printf
	// Log  = func(args ...interface{}) {}
	// Logf = func(args ...interface{}) {}
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// Params support named-paaram
type Params = map[string](interface{})

var DebugMode = true

// WorkingDir for agenda.
func WorkingDir() string {
	location, existed := os.LookupEnv("HOME")
	if !existed || DebugMode {
		location = "."
	}
	// NOTE: here to ensure workingdir existed ?
	ret := location + "/.agenda.d/"
	return ret
}
