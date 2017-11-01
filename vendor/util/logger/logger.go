package logger

import (
	"log"
)

// var Logger *log.Logger

// func init() {
// 	flog, err := os.Open(config.AgendaConfigPath())
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	Logger = log.New(flog, "logger: ", log.Llongfile)
// }

// var (
// 	Print   = Logger.Print
// 	Printf  = Logger.Printf
// 	Println = Logger.Println
// 	Fatal   = Logger.Fatal
// 	Fatalf  = Logger.Fatalf
// 	Fatalln = Logger.Fatalln
// 	Panic   = Logger.Panic
// 	Panicf  = Logger.Panicf
// 	Panicln = Logger.Panicln
//  // TODO:
//  Warning   = Logger.Print
//  Warningf  = Logger.Printf
//  Warningln = Logger.Println
//  Error     = Logger.Print
//  Errorf    = Logger.Printf
//  Errorln   = Logger.Println
// )

var (
	Print   = log.Print
	Printf  = log.Printf
	Println = log.Println
	Fatal   = log.Fatal
	Fatalf  = log.Fatalf
	Fatalln = log.Fatalln
	Panic   = log.Panic
	Panicf  = log.Panicf
	Panicln = log.Panicln

	// TODO:
	Warning   = log.Print
	Warningf  = log.Printf
	Warningln = log.Println
	Error     = log.Print
	Errorf    = log.Printf
	Errorln   = log.Println
)
