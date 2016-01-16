package nickserv

import (
	"io/ioutil"
	"github.com/Nyubis/mibot/ircmessage"
	"github.com/Nyubis/mibot/core"
)

var bot *core.Bot

func Init(ircbot *core.Bot) {
	bot = ircbot
}

func Handle(_ ircmessage.Message) {
	buf, _ := ioutil.ReadFile("./modules/nickserv/password.txt")
	bot.SendMessage("NickServ", "IDENTIFY " + string(buf)) 
}
