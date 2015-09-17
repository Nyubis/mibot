package channels

import (
	"github.com/nyubis/mibot/ircmessage"
	"github.com/nyubis/mibot/core"
	"github.com/nyubis/mibot/utils"
	"github.com/nyubis/mibot/modules/admin"
)

var autojoin []string
var blacklist []string
var bot *core.Bot

func Init(ircbot *core.Bot) {
	bot = ircbot

	LoadCfg()
}

func LoadCfg() {
	autojoin = prefix(core.Config.Channels.Autojoin)
	blacklist = prefix(core.Config.Channels.Blacklist)
}

func Autojoin(msg ircmessage.Message) {
	for _, channel := range autojoin {
		bot.SendJoin(channel)
	}
}

func InviteJoin(msg ircmessage.Message) {
	if len(msg.Content) > 0 && verify_channel(msg.Content) {
		bot.SendJoin(msg.Content)
	}
}

func HandleJoinCommand(channels []string, sender string, fromchannel string) {
	if admin.CheckAdmin(sender) {
		for _, channel := range channels {
			if verify_channel(channel) {
				bot.SendJoin(channel)
			}
		}
	}
}

func verify_channel(channel string) bool {
	return channel[0] == '#' && !utils.Contains(blacklist, channel)
}	

func prefix(channels []string) []string {
	for i, c := range channels {
		channels[i] = "#" + c
	}
	return channels
}
