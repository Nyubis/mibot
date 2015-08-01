package channels

import (
	"github.com/nyubis/mibot/ircmessage"
	"github.com/nyubis/mibot/core"
)

var autojoin = []string{"#bots", "#/g/sicp"}
var blacklist = []string{"#services", "#ripyourbot"}
var bot *core.Bot

func Init(ircbot *core.Bot) {
	bot = ircbot
}

func Autojoin(msg ircmessage.Message) {
	for _, channel := range autojoin {
		bot.SendJoin(channel)
	}
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
