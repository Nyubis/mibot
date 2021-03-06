package ignore

import (
	"github.com/Nyubis/mibot/core"
)

var bot *core.Bot

func Init(ircbot *core.Bot) {
	bot = ircbot
	LoadCfg()
}

func LoadCfg() {
	core.Ignored = make(map[string]bool)
	for _, nick := range core.Config.Ignored {
		core.Ignored[nick] = true
	}
}
