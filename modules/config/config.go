package config

import (
	"github.com/Nyubis/mibot/core"
	"github.com/Nyubis/mibot/modules/admin"
)

var ReloadFunc func()

func HandleReload(_ []string, sender string, channel string) {
	if admin.CheckAdmin(sender) {
		core.LoadConfig()
		ReloadFunc()
	}
}
