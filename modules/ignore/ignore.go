package ignore

import (
	"github.com/nyubis/mibot/core"
)

var bot *core.Bot

func Init(ircbot *core.Bot) {
	bot = ircbot
	for _, nick := range core.Config.Ignore.Ignored {
		core.Ignored[nick] = true
	}
}
