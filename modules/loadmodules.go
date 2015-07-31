package modules

import (
	"github.com/nyubis/mibot/modules/linktitle"
	"github.com/nyubis/mibot/modules/replies"
	"github.com/nyubis/mibot/modules/nickserv"
)

func Load() {
	AddModule("PRIVMSG", replies.Handle)
	AddModule("PRIVMSG", linktitle.Handle)
	AddModule("005", nickserv.Handle)
}
