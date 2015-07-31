package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"regexp"
	"strings"
	"github.com/nyubis/mibot/ircmessage"
	"github.com/nyubis/mibot/modules"
)

type Bot struct {
	nick          string
	user          string
	channel       string
	server        string
	port          int
	conn          net.Conn
	cinput        chan string
	coutput       chan string
}

func NewBot(nick string, channel string, server string, port int) *Bot {
	return &Bot{
		nick:          nick,
		user:          nick,
		channel:       channel,
		server:        server,
		port:          port,
	}
}

func (bot *Bot) SendCommand(cmd string) {
	output := sanitise(cmd) + "\r\n"
	fmt.Printf(">> %s", output)
	bot.coutput <- output
}

func (bot *Bot) Connect() {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", bot.server, bot.port))
	if err != nil {
		log.Fatal("Could not connect to server: ", err)
	}
	fmt.Printf("Connected to server %s:%d (%s)\n", bot.server, bot.port, conn.RemoteAddr())
	bot.coutput = make(chan string)
	bot.cinput = make(chan string)
	bot.conn = conn
}

func (bot *Bot) disconnect() {
	if bot.conn != nil {
		bot.conn.Close()
	}
}

func (bot *Bot) handleirc(line string) bool {
	// pingpong?
	if line[:4] == "PING" {
		bot.SendCommand("PONG" + line[4:])
		return true
	}
	// MOTD -> auth and join default channel
	type_re := regexp.MustCompile(` \d+ `)
	if type_re.FindString(line) == " 376 " {
		bot.SendCommand("JOIN " + bot.channel)
		return true
	}
	// Join channels on invite
	inv_re := regexp.MustCompile(` INVITE `)
	if inv_re.FindString(line) != "" {
		chan_re := regexp.MustCompile(":#.+$")
		channel := chan_re.FindString(line)[1:]
		bot.SendCommand("JOIN " + channel)
		return true
	}
	
	return false
}

func (bot *Bot) handle(line string) {
	replies := modules.Handle(ircmessage.Parse(line))
	for _, reply := range replies {
		bot.SendCommand(reply)
	}
}

func sanitise(cmd string) string {
	maxlength := 510
	cmd = strings.Replace(cmd, "\r", "", -1)
	cmd = strings.Replace(cmd, "\n", "", -1)
	if len(cmd) > maxlength {
		return cmd[:maxlength]
	}
	return cmd
}

func (bot *Bot) outputHandler() {
	for {
		out := <-bot.coutput
		fmt.Fprintf(bot.conn, out)
	}
}

func (bot *Bot) inputHandler() {
	for {
		line := <-bot.cinput
		fmt.Printf("<< %s\n", line)
		// Do irc stuff such as replying to pings, joining after motd, ...
		if bot.handleirc(line) {
			continue
		}
		// Otherwise, handle the channel data in a new goroutine
		go bot.handle(line)
	}
}


func main() {
	ircbot := NewBot("linkreader", "#/g/sicp", "irc.rizon.net", 6667)
	defer ircbot.disconnect()
	ircbot.Connect()
	modules.Load()

	reader := bufio.NewReader(ircbot.conn)
	tp := textproto.NewReader(reader)

	go ircbot.outputHandler()
	go ircbot.inputHandler()

	ircbot.SendCommand("USER " + ircbot.user + " 8 * :" + ircbot.nick)
	ircbot.SendCommand("NICK " + ircbot.nick)

	for {
		line, err := tp.ReadLine()
		if err != nil {
			log.Fatal(err)
			break
		}
		ircbot.cinput <- line
	}
}
