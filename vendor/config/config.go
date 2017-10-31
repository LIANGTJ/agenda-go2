package config

import (
	"convention/codec"
	"log"
	"os"
)

var DebugMode = true

// type Config = map[string](interface{})

// Config holds all configure of Agenda system.
var Config = make(map[string](interface{}))

// TODO: let Config, AllUsers, etc global singleton
func LoadConfig(decoder codec.Decoder) {
	cfg := &(Config)
	// CHECK: Need check if have already exactly loaded ALL config (i.e. eof) ?
	if err := decoder.Decode(cfg); err != nil {
		log.Fatal(err)
	}
}

func SaveConfig(encoder codec.Encoder) error {
	return encoder.Encode(Config)
}

// ... paths

// WorkingDir for agenda.
func WorkingDir() string {
	location, existed := os.LookupEnv("HOME")
	if !existed || DebugMode {
		location = "."
	}
	// NOTE: here to ensure workingdir existed ?
	ret := location + "/.agenda.d/"
	return ret
}

func UserDataRegisteredPath() string { return WorkingDir() + "user-registered.json" }
func UserDataPath() string           { return WorkingDir() + "user-data.json" }
func UserTestPath() string           { return WorkingDir() + "user-test.json" }

func MeetingDataPath() string { return WorkingDir() + "meeting-data.json" }
func MeetingTestPath() string { return WorkingDir() + "meeting-test.json" }

func AgendaConfigPath() string { return WorkingDir() + "config.json" }
