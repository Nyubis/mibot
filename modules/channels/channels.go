package channels

import (
	"github.com/Nyubis/mibot/ircmessage"
	"github.com/Nyubis/mibot/core"
	"github.com/Nyubis/mibot/utils"
	"github.com/Nyubis/mibot/modules/admin"
	"github.com/Nyubis/mibot/modules/floodcontrol"
)

var autojoin []string
var blacklist []string
var bot *core.Bot

func Init(ircbot *core.Bot) {
	bot = ircbot

	LoadCfg()
}

func LoadCfg() {
	autojoin = core.Config.Channels.Autojoin
	blacklist = core.Config.Channels.Blacklist
}

func Autojoin(msg ircmessage.Message) {
	for _, channel := range autojoin {
		bot.SendJoin(channel)
	}
}

func InviteJoin(msg ircmessage.Message) {
	if len(msg.Content) > 0 && verify_channel(msg.Content) {
		if !floodcontrol.FloodCheck("invite", msg.Nick, msg.Channel) {
			bot.SendJoin(msg.Content)
		}
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

func HandlePartCommand(channels []string, sender string, fromchannel string) {
	if admin.CheckAdmin(sender) {
		for _, channel := range channels {
			if verify_channel(channel) {
				bot.SendPart(channel)
			}
		}
	}
}

func verify_channel(channel string) bool {
	return channel[0] == '#' && !utils.Contains(blacklist, channel)
}
