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
func UserTestPath() string { return util.WorkingDir() + "user-test.json" }

func MeetingDataPath() string { return util.WorkingDir() + "meeting-data.json" }
func MeetingTestPath() string { return util.WorkingDir() + "meeting-test.json" }

// NOTE: REPEAT: `entity.DeserializeUserList`
func Load(ul *entity.UserList) error {
	// CHECK: Need clear ul ?
	file, err := os.Open(UserDataPath())
	if err != nil {
		log.Fatal(err)
		return err
	}
	decoder := json.NewDecoder(file)
	for decoder.More() {
		uInfo := new(entity.UserInfo)
		if err := decoder.Decode(uInfo); err != nil {
			log.Fatal(err)
			return err
		}
		user := entity.NewUser(*uInfo)
		if err := ul.Add(user); err != nil {
			log.Fatal(err)
			return err
		}
	}
	return nil
}

func SaveUserList(ul *entity.UserList) error {
	// FIXME: Not simply rewrite
	file, err := os.Create(UserDataPath())
	if err != nil {
		log.Fatal(err)
		return err
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
