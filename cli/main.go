package main

import (
	"agenda"
	"os"
	log "util/logger"

	cmd "github.com/Binly42/agenda-go/cmd"
)

// var logln = util.Log
// var logf = util.Logf

func init() {
}

func main() {
	agenda.LoadAll()
	defer agenda.SaveAll()

	if err := cmd.RootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1) // FIXME:
	}
}
