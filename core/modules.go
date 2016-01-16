package core

import (
	"fmt"
	"github.com/Nyubis/mibot/ircmessage"
	"strings"
)

const (
	commandChar = '@'
)

type module struct {
	triggerType string
	handler     func(ircmessage.Message)
}

var modules []module = make([]module, 0)
var commands map[string]func([]string, string, string) = make(map[string]func([]string, string, string), 0)
var Ignored map[string]bool = make(map[string]bool)

func AddModule(triggerType string, handler func(ircmessage.Message)) {
	modules = append(modules, module{triggerType, handler})
	fmt.Println("Added module for", triggerType)
}

func AddCommand(cmd string, handler func([]string, string, string)) {
	commands[cmd] = handler
	fmt.Println("Added command", cmd)
}

func PassToModules(msg ircmessage.Message) {
	if Ignored[msg.Nick] {
		return
	}
	if msg.Command == "PRIVMSG" && msg.Content[0] == commandChar {
		doCommand(msg)
	}
	for _, mod := range modules {
		if mod.triggerType == msg.Command {
			mod.handler(msg)
		}
	}
}

func doCommand(msg ircmessage.Message) {
	split := strings.Split(msg.Content, " ")
	cmdstring := split[0][1:]
	handler, exists := commands[cmdstring]

	if exists {
		handler(split[1:], msg.Nick, msg.Channel)
	}
}
