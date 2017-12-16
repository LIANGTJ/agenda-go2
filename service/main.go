package main

import (
	"agenda"
	"flag"
	"os"
	"os/signal"
	"syscall"
	log "util/logger"
)

// var logln = util.Log
// var logf = util.Logf

const (
	DefaultPort = agenda.DefaultPort
)

var (
	port string
	// ...
)

func init() {
	flag.StringVar(&port, "p", DefaultPort, "The PORT to be listened by agenda.")
}

func main() {
	flag.Parse()
	// TODO: validate port ?

	agenda.LoadAll()
	defer agenda.SaveAll()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		log.Infof("Signal %v", <-c)
		agenda.SaveAll()
		os.Exit(0)
	}()

	err := agenda.Listen(":" + port)
	if err != nil {
		log.Fatal(err)
	}
}
