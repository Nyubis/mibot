package modules

import (
	"github.com/nyubis/mibot/core"
	"github.com/nyubis/mibot/modules/linktitle"
	"github.com/nyubis/mibot/modules/replies"
	"github.com/nyubis/mibot/modules/nickserv"
	"github.com/nyubis/mibot/modules/channels"
)

func Load(bot *core.Bot) {
	core.AddModule("PRIVMSG", replies.Handle)
	replies.Init(bot)
	core.AddModule("PRIVMSG", linktitle.Handle)
	linktitle.Init(bot)
	core.AddModule("001", nickserv.Handle)
	nickserv.Init(bot)
	core.AddModule("001", channels.Autojoin)
	core.AddModule("INVITE", channels.InviteJoin)
	channels.Init(bot)
}
