package main

import (
	"os"
	log "util/logger"
	"cmd"
	"model"
)

// var logln = util.Log
// var logf = util.Logf

func init() {
}

func main() {
	defer model.SaveAll()
	model.LoadAll()
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1) // FIXME:
	}
}
