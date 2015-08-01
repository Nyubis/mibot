package channels

import (
	"github.com/nyubis/mibot/ircmessage"
	"github.com/nyubis/mibot/core"
)

var autojoin = "#bots"
var blacklist = []string{"#services", "#ripyourbot"}
var bot *core.Bot

func Init(ircbot *core.Bot) {
	bot = ircbot
}

func Autojoin(msg ircmessage.Message) {
	bot.SendJoin(autojoin)
}

func InviteJoin(msg ircmessage.Message) {
	if len(msg.Content) > 0 && msg.Content[0] == '#' && !contains(blacklist, msg.Content) {
		bot.SendJoin(msg.Content)
	}
}

func contains(hay []string, needle string) bool {
	for _, s := range hay {
		if s == needle {
			return true
		}
	}
	return false
}
