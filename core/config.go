package core

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
)

const (
	filename = "config.toml"
)

type ConfStruct struct {
	Nick     string
	Server   string
	Port     int
	TLS      bool
	Channels struct {
		Autojoin  []string
		Blacklist []string
	}
	Admins       []string
	Ignored      []string
	Chan         map[string]Channelsetting
	FloodControl struct {
		Invite FloodSetting
		Link   FloodSetting
		Reply  FloodSetting
	}
}

type FloodSetting struct {
	Time int
	Max  int
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
