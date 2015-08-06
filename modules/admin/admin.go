package admin

import (
	"github.com/nyubis/mibot/core"
	"github.com/nyubis/mibot/utils"
)

func CheckAdmin(nick string) bool {
	if utils.Contains(core.Config.Admins.Admins, nick) {
		return IsAuthenticated(nick)
	}
	return false
}
