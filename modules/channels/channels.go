package channels

import (
	"github.com/nyubis/mibot/ircmessage"
)

var autojoin = "#bots"
var blacklist = []string{"#services", "#ripyourbot"}

func Autojoin(msg ircmessage.Message) string {
	return ircmessage.Join(autojoin)
}

func InviteJoin(msg ircmessage.Message) string {
	if len(msg.Content) > 0 && msg.Content[0] == '#' && !contains(blacklist, msg.Content) {
		return ircmessage.Join(msg.Content)
	}
	return ""
}

func contains(hay []string, needle string) bool {
	for _, s := range hay {
		if s == needle {
			return true
		}
	}
	return false
}
