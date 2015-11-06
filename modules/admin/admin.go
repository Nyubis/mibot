package admin

import (
	"fmt"
	"github.com/nyubis/mibot/core"
	"github.com/nyubis/mibot/utils"
)

func CheckAdmin(nick string) bool {
	if utils.Contains(core.Config.Admins, nick) {
		fmt.Println(nick, "verified as admin")
		return IsAuthenticated(nick)
	}
	return false
}
