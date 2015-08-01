package replies

import (
	"github.com/nyubis/mibot/ircmessage"
	"github.com/nyubis/mibot/core"
)

var bot *core.Bot

func Init(ircbot *core.Bot) {
	bot = ircbot
}

func Handle(msg ircmessage.Message) {
	if msg.Content == ".bots" {
		bot.SendMessage(msg.Channel, "Reporting in! [Go] On Github soon, I promise!")
	}
}
