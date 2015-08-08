package config

import (
	"github.com/nyubis/mibot/core"
	"github.com/nyubis/mibot/modules/admin"
)

var ReloadFunc func()

func HandleReload(_ []string, sender string, channel string) {
	if admin.CheckAdmin(sender) {
		core.LoadConfig()
		ReloadFunc()
	}
}
		
