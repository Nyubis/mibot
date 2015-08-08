package replies

import (
	"github.com/nyubis/mibot/ircmessage"
	"github.com/nyubis/mibot/core"
	"strings"
	"fmt"
)

var bot *core.Bot
var replies map[string]string

func Init(ircbot *core.Bot) {
	bot = ircbot
	LoadCfg()
}

func LoadCfg() {
	replies = make(map[string]string)
	for _, line  := range core.Config.Replies.Replies {
		if !strings.Contains(line, "\t") {
			fmt.Println("Warning: No tab found in line %s, reply was not added", line)
			continue
		}
		split := strings.Split(line, "\t")
		replies[split[0]] = split[1]
	}
}


func Handle(msg ircmessage.Message) {
	reply, has := replies[msg.Content]
	if has {
		bot.SendMessage(msg.Channel, reply)
	}
}
