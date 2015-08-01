package modules

import (
	"fmt"
	"github.com/nyubis/mibot/ircmessage"
)

type module struct {
	triggerType string
	handler     func(ircmessage.Message) string
}

var modules []module = make([]module, 0)

func AddModule(triggerType string, handler func(ircmessage.Message) string) {
	modules = append(modules, module{triggerType, handler})
	fmt.Printf("Added module for %s\n", triggerType)
}

func Handle(msg ircmessage.Message) []string {
	replies := make([]string, 0)
	for _, mod := range modules {
		if mod.triggerType == msg.Command {
			reply := mod.handler(msg)
			if reply != "" {
				replies = append(replies, reply)
			}
		}
	}

	return replies
}

