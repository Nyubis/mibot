package replies

import (
	"github.com/Nyubis/mibot/core"
	"github.com/Nyubis/mibot/ircmessage"
	"github.com/Nyubis/mibot/modules/floodcontrol"
)

var bot *core.Bot
var globalreplies map[string]string
var channelreplies map[string]map[string]string

func Init(ircbot *core.Bot) {
	bot = ircbot
	LoadCfg()
}

func LoadCfg() {
	channelreplies = make(map[string]map[string]string)
	for k, v := range core.Config.Channelsettings {
		if k == "global" {
			globalreplies = v.Replies
		} else {
			channelreplies[k] = v.Replies
		}
	}
}

func Handle(msg ircmessage.Message) {
	has := false
	var reply string
	// Check whether we have a set of replies specific to this channel
	replyset, haschannel := channelreplies[msg.Channel]
	if haschannel {
		// Check whether the trigger is in the set
		reply, has = replyset[msg.Content]
	}

	// If nothing in the channel-specific set, check the global reply set.
	if !has {
		reply, has = globalreplies[msg.Content]
	}
	// If match, also check for flood: don't trigger replies too often
	if has && !floodcontrol.FloodCheck("reply", msg.Channel) {
		bot.SendMessage(msg.Channel, reply)
	}
}
