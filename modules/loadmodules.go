package modules

import (
	"github.com/Nyubis/mibot/core"
	"github.com/Nyubis/mibot/modules/admin"
	"github.com/Nyubis/mibot/modules/channels"
	"github.com/Nyubis/mibot/modules/config"
	"github.com/Nyubis/mibot/modules/floodcontrol"
	"github.com/Nyubis/mibot/modules/ignore"
	"github.com/Nyubis/mibot/modules/linktitle"
	"github.com/Nyubis/mibot/modules/nickserv"
	"github.com/Nyubis/mibot/modules/replies"
)

func Load(bot *core.Bot) {
	core.AddModule("PRIVMSG", replies.Handle)
	core.AddModule("PRIVMSG", linktitle.Handle)
	core.AddModule("352", admin.ReceiveWho)
	core.AddModule("001", nickserv.Handle)
	core.AddModule("001", channels.Autojoin)
	core.AddModule("353", channels.ReceiveNamesReply)
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
	floodcontrol.LoadCfg()
}
