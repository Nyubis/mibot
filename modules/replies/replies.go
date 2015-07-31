package replies

import (
	"github.com/nyubis/mibot/ircmessage"
)

func Handle(msg ircmessage.Message) string {
	if msg.Content == ".bots" {
		return ircmessage.PrivMsg(msg.Channel, "Reporting in! [Go] Code is proprietary :^)")
	}
	return ""
}
