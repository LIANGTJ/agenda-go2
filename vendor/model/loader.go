package model

import (
	"bufio"
	"encoding/json"
	"entity"
	"log"
	"os"
	"sync"
	"util"
)

var wg sync.WaitGroup
var fin = os.Stdin

type Loader interface {
	// TODO: need to abstract UserList and MeetingList first
	Load(*interface{}) error
}

type UserLoader struct {
}

func UserDataPath() string { return util.WorkingDir() + "user-data.json" }

func (loader *UserLoader) Load(ul *entity.UserList) error {
	file, err := os.Open(UserDataPath())
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(file)
	for decoder.More() {
		var m entity.Meeting
		if err := decoder.Decode(&m); err != nil {
			log.Fatal(err)
		}
		// fmt.Printf("%s: %s\n", m.Name, m.Text)
	}
	return nil
}

func Save(ul *entity.UserList) error {
	file, err := os.Create(UserDataPath())
	if err != nil {
		log.Fatal(err)
	}
	encoder := json.NewEncoder(file)
	if err := ul.Serialize(encoder); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

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
