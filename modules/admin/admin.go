package admin

import (
	"fmt"
	"github.com/Nyubis/mibot/core"
	"github.com/Nyubis/mibot/utils"
)

func CheckAdmin(nick string) bool {
	if utils.Contains(core.Config.Admins, nick) {
		fmt.Println(nick, "verified as admin")
		return IsAuthenticated(nick)
	}
	return false
}
