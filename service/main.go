package main

import (
	"agenda"
	"flag"
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

	server := agenda.New()
	err := server.Listen(":" + port)
	if err != nil {
		log.Fatal(err)
	}
}
