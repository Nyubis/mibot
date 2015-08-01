package main

import (
	"github.com/nyubis/mibot/core"
	"github.com/nyubis/mibot/modules"
)


func main() {
	ircbot := core.NewBot("binkreader", "#bots", "irc.rizon.net", 6667)
	defer ircbot.Disconnect()
	modules.Load(ircbot)
	ircbot.Connect()
}
