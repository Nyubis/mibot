package modules

import (
	"github.com/nyubis/mibot/core"
	"github.com/nyubis/mibot/modules/linktitle"
	"github.com/nyubis/mibot/modules/replies"
	"github.com/nyubis/mibot/modules/nickserv"
	"github.com/nyubis/mibot/modules/channels"
	"github.com/nyubis/mibot/modules/admin"
	"github.com/nyubis/mibot/modules/ignore"
	"github.com/nyubis/mibot/modules/config"
	"github.com/nyubis/mibot/modules/floodcontrol"
)

func Load(bot *core.Bot) {
	core.AddModule("PRIVMSG", replies.Handle)
	core.AddModule("PRIVMSG", linktitle.Handle)
	core.AddModule("352", admin.ReceiveWho)
	core.AddModule("001", nickserv.Handle)
	core.AddModule("001", channels.Autojoin)
	core.AddModule("INVITE", channels.InviteJoin)

	core.AddCommand("shorten", linktitle.HandleShorten)
	core.AddCommand("enable", linktitle.HandleEnable)
	core.AddCommand("disable", linktitle.HandleDisable)
	core.AddCommand("reload", config.HandleReload)
	core.AddCommand("join", channels.HandleJoinCommand)
	core.AddCommand("part", channels.HandlePartCommand)

	replies.Init(bot)
	linktitle.Init(bot)
	admin.Init(bot)
	nickserv.Init(bot)
	channels.Init(bot)
	ignore.Init(bot)
	floodcontrol.Init(bot)
	config.ReloadFunc = reload
}

func reload() {
	ignore.LoadCfg()
	channels.LoadCfg()
	replies.LoadCfg()
	linktitle.LoadCfg()
}

