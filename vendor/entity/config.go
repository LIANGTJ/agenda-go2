package entity

import (
	"log"
)

// type Config = map[string](interface{})

// func NewConfig() *Config {
// 	cfg := make(Config)
// 	return &cfg
// }

// Config holds all configure of Agenda system.
var Config = make(map[string](interface{}))

// TODO: let Config, AllUsers, etc global singleton
func LoadConfig(decoder Decoder) {
	cfg := &(Config)
	// CHECK: Need check if have already exactly loaded ALL config (i.e. eof) ?
	if err := decoder.Decode(cfg); err != nil {
		log.Fatal(err)
	}
}

func SaveConfig(encoder Encoder) error {
	return encoder.Encode(Config)
}
