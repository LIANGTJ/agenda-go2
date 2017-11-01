package logger

import (
	"config"
	"io"
	"log"
	"os"
)

var Logger *log.Logger

func init() {
	flog, err := os.OpenFile(config.LogPath(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		log.Panic(err)
	}
	Logger = log.New(io.MultiWriter(flog, os.Stderr), "agenda: ", log.LstdFlags|log.Lshortfile)
}

func Print(v ...interface{})                 { Logger.SetPrefix("[info]"); Logger.Print(v...) }
func Printf(format string, v ...interface{}) { Logger.SetPrefix("[info]"); Logger.Printf(format, v...) }
func Println(v ...interface{})               { Logger.SetPrefix("[info]"); Logger.Println(v...) }
func Fatal(v ...interface{})                 { Logger.SetPrefix("[fatal]"); Logger.Fatal(v...) }
func Fatalf(format string, v ...interface{}) { Logger.SetPrefix("[fatal]"); Logger.Fatalf(format, v...) }
func Fatalln(v ...interface{})               { Logger.SetPrefix("[fatal]"); Logger.Fatalln(v...) }
func Panic(v ...interface{})                 { Logger.SetPrefix("[panic]"); Logger.Panic(v...) }
func Panicf(format string, v ...interface{}) { Logger.SetPrefix("[panic]"); Logger.Panicf(format, v...) }
func Panicln(v ...interface{})               { Logger.SetPrefix("[panic]"); Logger.Panicln(v...) }

func Warning(v ...interface{}) { Logger.SetPrefix("[warn]"); Logger.Print(v...) }
func Warningf(format string, v ...interface{}) {
	Logger.SetPrefix("[warn]")
	Logger.Printf(format, v...)
}
func Warningln(v ...interface{})             { Logger.SetPrefix("[warn]"); Logger.Println(v...) }
func Error(v ...interface{})                 { Logger.SetPrefix("[error]"); Logger.Print(v...) }
func Errorf(format string, v ...interface{}) { Logger.SetPrefix("[error]"); Logger.Printf(format, v...) }
func Errorln(v ...interface{})               { Logger.SetPrefix("[error]"); Logger.Println(v...) }

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
// 	// TODO:
// 	Warning   = Logger.Print
// 	Warningf  = Logger.Printf
// 	Warningln = Logger.Println
// 	Error     = Logger.Print
// 	Errorf    = Logger.Printf
// 	Errorln   = Logger.Println
// )

// var (
// 	Print   = log.Print
// 	Printf  = log.Printf
// 	Println = log.Println
// 	Fatal   = log.Fatal
// 	Fatalf  = log.Fatalf
// 	Fatalln = log.Fatalln
// 	Panic   = log.Panic
// 	Panicf  = log.Panicf
// 	Panicln = log.Panicln

// 	// TODO:
// 	Warning   = log.Print
// 	Warningf  = log.Printf
// 	Warningln = log.Println
// 	Error     = log.Print
// 	Errorf    = log.Printf
// 	Errorln   = log.Println
// )
