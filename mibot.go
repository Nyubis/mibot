package main

import (
	"github.com/Nyubis/mibot/core"
	"github.com/Nyubis/mibot/modules"
)

func main() {
	core.LoadConfig()
	ircbot := core.NewBot(
		core.Config.Nick,
		core.Config.Server,
		core.Config.Port,
		core.Config.TLS)
	defer ircbot.Disconnect()
	modules.Load(ircbot)
	ircbot.Start()
}
