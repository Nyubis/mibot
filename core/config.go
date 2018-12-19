package core

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
)

const (
	filename = "config.toml"
)

type ConfStruct struct {
	Nick     string
	Server   string
	Port     int
	Channels struct {
		Autojoin  []string
		Blacklist []string
	}
	Admins  []string
	Ignored []string
	Chan    map[string]Channelsetting
}

type Channelsetting struct {
	Disabled  bool
	Replies   [][]string
	Blacklist []string
}

var Config ConfStruct

func LoadConfig() {
	_, err := toml.DecodeFile(filename, &Config)
	if err != nil {
		log.Fatal("Failed to read config:\n", err)
	}
	fmt.Println(Config)
}
