package main

import (
	"github.com/nyubis/mibot/core"
	"github.com/nyubis/mibot/modules"
)

func main() {
	core.LoadConfig()
	ircbot := core.NewBot(
		core.Config.Nick,
		core.Config.Server,
		core.Config.Port)
	defer ircbot.Disconnect()
	modules.Load(ircbot)
	ircbot.Connect()
}
