package main

import (
	"entity"
	"fmt"
	"log"
	"model"
	"os"

	cmd "github.com/Binly42/agenda-go/cmd"
	// "github.com/spf13/cobra"
)

// var logln = util.Log
// var logf = util.Logf

func init() {
	// flags ?
}

func main() {
	// p := packer.NewOne()
	// m := entity.Meeting
	// agenda.work()

	model.Load()

	log.Printf("Users:  %+v\n", entity.GetAllUsersRegistered())
	log.Printf("Config: %+v\n", entity.Config)

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
