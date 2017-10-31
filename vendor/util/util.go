package util

import (
	"log"

	Logger "util/logger"
)

var (
	// Logger = logger.Logger

	Print   = Logger.Print
	Printf  = Logger.Printf
	Println = Logger.Println
	Fatal   = Logger.Fatal
	Fatalf  = Logger.Fatalf
	Fatalln = Logger.Fatalln
	Panic   = Logger.Panic
	Panicf  = Logger.Panicf
	Panicln = Logger.Panicln
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

var (
	Log  = log.Println
	Logf = log.Printf

// Log  = func(args ...interface{}) {}
// Logf = func(args ...interface{}) {}
)

// Params support named-paaram
type Params = map[string](interface{})
