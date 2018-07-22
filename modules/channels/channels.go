package channels

import (
	"github.com/Nyubis/mibot/core"
	"github.com/Nyubis/mibot/ircmessage"
	"github.com/Nyubis/mibot/modules/admin"
	"github.com/Nyubis/mibot/modules/floodcontrol"
	"github.com/Nyubis/mibot/utils"
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
		// If the NAMES reply contains only our nick, we will part it again
	}
}

func InviteJoin(msg ircmessage.Message) {
	if len(msg.Content) > 0 && verify_channel(msg.Content) {
		if !floodcontrol.FloodCheck("invite", msg.Channel) {
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

func ReceiveNamesReply(msg ircmessage.Message) {
	// Leave the channel if we're the only user (and it was autojoined)
	// We may have received an @ or other operator char, hence the second check
	if msg.Content == core.Config.Nick || msg.Content[1:] == core.Config.Nick {
		chanName := msg.Params[len(msg.Params) -1]
		if utils.Contains(autojoin, chanName) {
			bot.SendPart(chanName)
		}
	}
}
