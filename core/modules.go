package core

import (
	"fmt"
	"github.com/nyubis/mibot/ircmessage"

)

type module struct {
	triggerType string
	handler     func(ircmessage.Message)
}

var modules []module = make([]module, 0)

func AddModule(triggerType string, handler func(ircmessage.Message)) {
	modules = append(modules, module{triggerType, handler})
	fmt.Printf("Added module for %s\n", triggerType)
}

func PassToModules(msg ircmessage.Message) {
	for _, mod := range modules {
		if mod.triggerType == msg.Command {
			mod.handler(msg)
		}
	}
}

