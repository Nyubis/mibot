package core

import (
	"encoding/json"
	"fmt"
	"github.com/DisposaBoy/JsonConfigReader"
	"log"
	"os"
)

const (
	filename = "config.json"
)

type ConfStruct struct {
	Nick     string
	Server   string
	Port     int
	Channels struct {
		Autojoin  []string
		Blacklist []string
	}
	Admins          []string
	Ignored         []string
	Channelsettings map[string]Channelsetting
}

type Channelsetting struct {
	Disabled  bool
	Replies   map[string]string
	Blacklist []string
}

var Config ConfStruct

func LoadConfig() {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Failed to read config:\n", err)
	}
	reader := JsonConfigReader.New(file)
	json.NewDecoder(reader).Decode(&Config)
	fmt.Println(Config)
}
