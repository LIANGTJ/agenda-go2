package logger

import (
	"config"
	"fmt"
	"log"
	"os"
)

var Logger *log.Logger

func init() {
	flog, err := os.OpenFile(config.LogPath(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		log.Panic(err)
	}
	// logWriter := io.MultiWriter(flog, os.Stderr)
	logWriter := flog
	Logger = log.New(logWriter, "agenda: ", log.LstdFlags|log.Lshortfile)
}

func Print(v ...interface{})                 { Logger.SetPrefix("[info]"); Logger.Print(v...) }
func Printf(format string, v ...interface{}) { Logger.SetPrefix("[info]"); Logger.Printf(format, v...) }
func Println(v ...interface{})               { Logger.SetPrefix("[info]"); Logger.Println(v...) }

func Warning(v ...interface{}) { Logger.SetPrefix("[warning]"); Logger.Print(v...) }
func Warningf(format string, v ...interface{}) {
	Logger.SetPrefix("[warning]")
	Logger.Printf(format, v...)
}
func Warningln(v ...interface{}) { Logger.SetPrefix("[warning]"); Logger.Println(v...) }

func Error(v ...interface{}) {
	Logger.SetPrefix("[error]")
	Logger.Print(v...)
	fmt.Fprint(os.Stderr, v...)
}
func Errorf(format string, v ...interface{}) {
	Logger.SetPrefix("[error]")
	Logger.Printf(format, v...)
	fmt.Fprintf(os.Stderr, format, v...)
}
func Errorln(v ...interface{}) {
	Logger.SetPrefix("[error]")
	Logger.Println(v...)
	fmt.Fprintln(os.Stderr, v...)
}

func Fatal(v ...interface{}) {
	Logger.SetPrefix("[fatal]")
	Logger.Fatal(v...)
	fmt.Fprint(os.Stderr, v...)
}
func Fatalf(format string, v ...interface{}) {
	Logger.SetPrefix("[fatal]")
	Logger.Fatalf(format, v...)
	fmt.Fprintf(os.Stderr, format, v...)
}
func Fatalln(v ...interface{}) {
	Logger.SetPrefix("[fatal]")
	Logger.Fatalln(v...)
	fmt.Fprintln(os.Stderr, v...)
}

func Panic(v ...interface{}) {
	Logger.SetPrefix("[panic]")
	Logger.Panic(v...)
	fmt.Fprint(os.Stderr, v...)
}
func Panicf(format string, v ...interface{}) {
	Logger.SetPrefix("[panic]")
	Logger.Panicf(format, v...)
	fmt.Fprintf(os.Stderr, format, v...)
}
func Panicln(v ...interface{}) {
	Logger.SetPrefix("[panic]")
	Logger.Panicln(v...)
	fmt.Fprintln(os.Stderr, v...)
}

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
