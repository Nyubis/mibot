package core

import (
	"github.com/scalingdata/gcfg"
	"log"
)

const (
	filename = "mibot.cfg"
)

type ConfStruct struct {
	Main struct {
		Nick string
		Server string
		Port int
	}
	Nickserv struct {
		Password_file string
	}
	Channels struct {
		Autojoin []string
		Blacklist []string
	}
	Replies struct {
		Replies []string
	}
	Admins struct {
		Admins []string
	}
	Ignore struct {
		Ignored []string
	}
}

var Config ConfStruct

func LoadConfig() {
	err := gcfg.ReadFileInto(&Config, filename)
	if err != nil {
		log.Fatal("Failed to read config:\n", err)
	}
}

