package ircmessage

import (
	"fmt"
	"regexp"
	"strings"
)

type Message struct {
	Sender  string
	Nick    string
	Command string
	Params  []string
	Channel string
	Content string
}

var re = regexp.MustCompile(":([^ ]+) ([^ ]+) ([^:]*)(:.*)?")

func Parse(line string) Message {
	var msg Message
	matches := re.FindStringSubmatch(line)

	msg.Sender = matches[1]
	msg.Command = matches[2]
	if matches[3] != "" {
		msg.Params = strings.Split(strings.Trim(matches[3], " "), " ")
	}
	if matches[4] != "" {
		msg.Content = matches[4][1:] // Remove leading ':'
	}
	if strings.Contains(msg.Sender, "!") {
		// Extract nick from full ident, if applicable
		msg.Nick = msg.Sender[:strings.Index(msg.Sender, "!")]
	}
	if len(msg.Params) > 0 && msg.Params[0][0] == '#' {
		// Extract channel name from params, if applicable
		msg.Channel = msg.Params[0]
	} else if msg.Nick != "" {
		// Use nick as channel, helps with replies
		msg.Channel = msg.Nick
	}

	return msg
}

func PrivMsg(to string, content string) string {
	return fmt.Sprintf("PRIVMSG %s :%s", to, content)
}

func Join(channel string) string {
	return fmt.Sprintf("JOIN %s", channel)
}
