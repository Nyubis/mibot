package main

import (
	"github.com/nyubis/mibot/core"
	"github.com/nyubis/mibot/modules"
)

func main() {
	core.LoadConfig()
	ircbot := core.NewBot(
		core.Config.Main.Nick,
		core.Config.Main.Server,
		core.Config.Main.Port)
	defer ircbot.Disconnect()
	modules.Load(ircbot)
	ircbot.Connect()
}
