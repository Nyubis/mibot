package nickserv

import (
	"io/ioutil"
	"github.com/nyubis/mibot/ircmessage"
)

func Handle(_ ircmessage.Message) string {
	buf, _ := ioutil.ReadFile("./modules/nickserv/password.txt")
	return ircmessage.PrivMsg("NickServ", "IDENTIFY " + string(buf)) 
}
