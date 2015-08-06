package modules

import (
	"github.com/nyubis/mibot/core"
	"github.com/nyubis/mibot/modules/linktitle"
	"github.com/nyubis/mibot/modules/replies"
	"github.com/nyubis/mibot/modules/nickserv"
	"github.com/nyubis/mibot/modules/channels"
	"github.com/nyubis/mibot/modules/admin"
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

	replies.Init(bot)
	linktitle.Init(bot)
	admin.Init(bot)
	nickserv.Init(bot)
	channels.Init(bot)
}
