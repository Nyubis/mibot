package admin

import (
	"github.com/nyubis/mibot/core"
)

func CheckAdmin(nick string) bool {
	if contains(core.Config.Admins.Admins, nick) {
		return IsAuthenticated(nick)
	}
	return false
}

func contains(list []string, item string) bool {
	for _, s := range list {
		if s == item {
			return true
		}
	}
	return false
}
