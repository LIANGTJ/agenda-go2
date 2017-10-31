package logger

import (
	"config"
	"log"
	"os"
)

var Logger *log.Logger

func init() {
	flog, err := os.Open(config.AgendaConfigPath())
	if err != nil {
		log.Panic(err)
	}
	Logger = log.New(flog, "logger: ", log.Llongfile)
}
