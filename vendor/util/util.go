package util

import (
	"log"
	// "util/logger"
)

// var (
// 	Logger = logger.Logger
// )

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
