package main

import (
	"entity"
	"model"
	"util"
)

var logln = util.Log
var logf = util.Logf

func main() {
	model.LoadConfig()
	// model.SaveConfig()

	logf("Config: %+v\n", entity.Config)
}
