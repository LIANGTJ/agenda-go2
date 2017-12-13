package main

import (
	"agenda"
	"os"
	log "util/logger"
	"fmt"
	cmd "github.com/Binly42/agenda-go/cmd"
)

// var logln = util.Log
// var logf = util.Logf

func init() {
}

func main() {
	agenda.LoadAll()
	defer agenda.SaveAll()
	fmt.Println("main called")
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1) // FIXME:
	}
}
