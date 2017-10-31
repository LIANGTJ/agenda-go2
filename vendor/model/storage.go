package model

import (
	"bufio"
	"config"
	"encoding/json"
	"entity"
	"log"
	"os"
	"sync"
	agendaLogger "util/logger"
)

var wg sync.WaitGroup
var fin = os.Stdin

// Load : load all resources for agenda.
func Load() {
	loadConfig()
	loadRegisteredUserList()
}

// Save : Save all data for agenda.
func Save() error {
	if err := saveRegisteredUserList(); err != nil {
		agendaLogger.Println(err.Error())
		return err
	}
	if err := saveAllMeeting(); err != nil {
		agendaLogger.Println(err.Error())
		return err
	}
	if err := saveConfig(); err != nil {
		agendaLogger.Println(err.Error())
		return err
	}
	return nil
}
func saveAllMeeting() error {
	log.Println("saveMeeting Error------------\n")
	fout, err := os.Create(config.MeetingDataPath())
	if err != nil {
		log.Fatal(err)
	}
	encoder := json.NewEncoder(fout)

	if err := entity.GetAllMeetings().Save(encoder); err != nil {
		log.Printf(err.Error()) // TODO: hadnle ?
		return err
	}
	return nil
}
func loadConfig() {
	fcfg, err := os.Open(config.AgendaConfigPath())
	if err != nil {
		log.Fatalf("Load config fail, for config path: %v\n", config.AgendaConfigPath())
	}
	decoder := json.NewDecoder(fcfg)

	config.LoadConfig(decoder)
}
func saveConfig() error {
	fcfg, err := os.Create(config.AgendaConfigPath())
	if err != nil {
		log.Fatalf("Save config fail, for config path: %v\n", config.AgendaConfigPath())
	}
	encoder := json.NewEncoder(fcfg)

	return config.SaveConfig(encoder)
}

func loadRegisteredUserList() {
	fin, err := os.Open(config.UserDataRegisteredPath())
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(fin)

	entity.LoadUsersAllRegistered(decoder)
}
func saveRegisteredUserList() error {
	fout, err := os.Create(config.UserDataRegisteredPath())
	if err != nil {
		log.Fatal(err)
	}
	encoder := json.NewEncoder(fout)

	if err := entity.SaveUsersAllRegistered(encoder); err != nil {
		log.Printf(err.Error()) // TODO: hadnle ?
		return err
	}
	return nil
}

// .....
// ref to before

func readInput() (<-chan string, error) {
	channel := make(chan string)
	scanner := bufio.NewScanner(fin)
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	wg.Add(1)
	go func() {
		for scanner.Scan() {
			channel <- scanner.Text() + "\n"
		}
		defer wg.Done()
		close(channel)
	}()

	return channel, nil
}
