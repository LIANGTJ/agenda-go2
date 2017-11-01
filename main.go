package main

import (
	"agenda"
	"config"
	"entity"
	"fmt"
	"util"

	cmd "github.com/Binly42/agenda-go/cmd"
	// "github.com/spf13/cobra"
)

var logln = util.Log
var logf = util.Logf

func init() {
}

func main() {
	agenda.LoadAll()
	defer agenda.SaveAll()

	// logf("Users:  %+v\n", entity.GetAllUsersRegistered())
	// logf("Config: %+v\n", config.Config)

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		// os.Exit(1) FIXME:
	}
}
