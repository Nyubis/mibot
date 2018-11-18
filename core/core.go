package core

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"strings"
	"time"

	"github.com/Nyubis/mibot/ircmessage"
)

const (
	timeout = time.Second * 360
)

type Bot struct {
	nick    string
	user    string
	server  string
	port    int
	conn    net.Conn
	cinput  chan string
	coutput chan string
}

func NewBot(nick string, server string, port int) *Bot {
	return &Bot{
		nick:   nick,
		user:   nick,
		server: server,
		port:   port,
	}
}

func (bot *Bot) SendCommand(cmd string) {
	output := sanitise(cmd) + "\r\n"
	fmt.Printf(">> %s", output)
	bot.coutput <- output
}

func (bot *Bot) SendMessage(to string, content string) {
	bot.SendCommand(fmt.Sprintf("PRIVMSG %s :%s", to, content))
}

func (bot *Bot) SendJoin(channel string) {
	bot.SendCommand("JOIN " + channel)
}

func (bot *Bot) SendPart(channel string) {
	bot.SendCommand("PART " + channel)
}

func (bot *Bot) Connect() {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", bot.server, bot.port), time.Second * 5)
	if err != nil {
		log.Fatal("Could not connect to server: ", err)
	}
	fmt.Printf("Connected to server %s:%d (%s)\n", bot.server, bot.port, conn.RemoteAddr())
	bot.coutput = make(chan string)
	bot.cinput = make(chan string)
	bot.conn = conn
	reader := bufio.NewReader(bot.conn)
	tp := textproto.NewReader(reader)

	go bot.outputHandler()
	go bot.inputHandler()

	bot.SendCommand("USER " + bot.user + " 8 * :" + bot.nick)
	bot.SendCommand("NICK " + bot.nick)

	for {
		line, err := tp.ReadLine()
		conn.SetReadDeadline(time.Now().Add(timeout))
		if err != nil {
			if strings.HasSuffix(err.Error(), "i/o timeout") {
				return
			}
			log.Fatal(err)
		}
		bot.cinput <- line
	}
}

func (bot *Bot) Disconnect() {
	if bot.conn != nil {
		bot.conn.Close()
	}
}

func (bot *Bot) Start() {
	for {
		// returns when connection is lost, e.g. ping timeout
		bot.Connect()
		bot.Disconnect()
		fmt.Println("Reconnecting…")
	}
}

func (bot *Bot) handleirc(line string) bool {
	// pingpong?
	if line[:4] == "PING" {
		bot.SendCommand("PONG" + line[4:])
		return true
	}

	return false
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
		fmt.Fprint(bot.conn, out)
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
		go PassToModules(ircmessage.Parse(line))
	}
}
