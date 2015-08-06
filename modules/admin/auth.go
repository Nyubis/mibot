package admin

import (
	"fmt"
	"strings"

	"github.com/nyubis/mibot/core"
	"github.com/nyubis/mibot/ircmessage"
)

var bot *core.Bot
var whochan chan ircmessage.Message

func Init(ircbot *core.Bot) {
	bot = ircbot
	whochan = make(chan ircmessage.Message)
}

func IsAuthenticated(nick string) bool {
	bot.SendCommand("WHO " + nick)
	fmt.Println("WHO sent, waiting for channel...")
	msg := <-whochan
	fmt.Println("WHO received from channel")

	l := len(msg.Params) - 1
	flags := msg.Params[l] // Last element
	recv_nick := msg.Params[l-1] // Second to last

	return strings.Contains(flags, "r")  && recv_nick == nick
}
	
func ReceiveWho(msg ircmessage.Message) {
	fmt.Println("Who receiving, queueing up")
	whochan <- msg
}
