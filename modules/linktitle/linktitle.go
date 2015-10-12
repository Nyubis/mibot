package linktitle

import (
	"regexp"
	"net/http"
	"io"
	"fmt"
	"strings"
	"errors"
	"time"

	"github.com/nyubis/mibot/ircmessage"
	"github.com/nyubis/mibot/core"
	"github.com/nyubis/mibot/modules/admin"
	"github.com/nyubis/mibot/modules/floodcontrol"
	"github.com/nyubis/mibot/utils"

	"golang.org/x/net/html"
)

const (
	redirectLimit = 10
	titleLimit    = 200
	timeout       = 5000
	byteLimit     = 65536

	shortenURL    = "http://is.gd/create.php?format=simple&url="
)

var httpRe = regexp.MustCompile("https?://[^\\s]*")

var lastURL string
var disabled map[string]bool

var bot *core.Bot
var client = &http.Client{
	CheckRedirect: func(_ *http.Request, via []*http.Request) error {
		if len(via) >= redirectLimit {
			return errors.New(fmt.Sprintf("Stopped after %d redirects", redirectLimit))
		}
		return nil
	},
	Timeout: time.Duration(timeout)*time.Millisecond,
}

func Init(ircbot *core.Bot) {
	bot = ircbot
	disabled = make(map[string]bool)
}

// Incoming event for and irc message
func Handle(msg ircmessage.Message) {
	url := httpRe.FindString(msg.Content)
	if url == "" {
		return
	}

	if floodcontrol.FloodCheck("link", msg.Nick, msg.Channel) {
		return
	}

	lastURL = url
	if !disabled[msg.Channel] {
		title := getAndFindTitle(url)
		if title != "" {
			bot.SendMessage(msg.Channel, "[URL] " + title)
		}
	} else {
		fmt.Println("Link detected, but", msg.Channel, "is disabled")
	}
}

// Incoming event for @shorten
func HandleShorten(_ []string, sender string, channel string) {
	if lastURL == "" {
		fmt.Println("No last url to shorten")
		return
	}
	short, err := shorten(lastURL)
	if err != nil {
		fmt.Println("Failed to shorten url", lastURL)
		fmt.Println(err)
		return
	}
	bot.SendMessage(channel, sender + ": " + short)
}

// Incoming event for @disable
func HandleDisable(_ []string, sender string, channel string) {
	if admin.CheckAdmin(sender) {
		disabled[channel] = true
		bot.SendMessage(channel, "Link reading disabled for this channel")
	}
}

// Incoming event for @enable
func HandleEnable(_ []string, sender string, channel string) {
	if admin.CheckAdmin(sender) {
		disabled[channel] = false
		bot.SendMessage(channel, "Link reading enabled for this channel")
	}
}

// Do the http request and continue to findTitle()
func getAndFindTitle(url string) string {
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()

	buf := make([]byte, byteLimit)
	n, _ := io.ReadFull(resp.Body, buf)

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		if resp.ContentLength >= 0 {
			return fmt.Sprintf("%s, %s", contentType, utils.ParseSize(resp.ContentLength))
		} else {
			return fmt.Sprintf("%s, unknown size", contentType)
		}
	}

	title := findTitle(string(buf[:n]))
	if title != "" {
		// Strip urls in the title, remove spaces on the beginning and end, and truncate it
		return utils.Truncate(strings.TrimSpace(stripUrls(title)), titleLimit)
	}
	return ""
}

// Parse the HTML in a string and extract the title
func findTitle(data string) string {
	tz := html.NewTokenizer(strings.NewReader(data))
	inbody := false
	for {
		t := tz.Next()
		tn, _ := tz.TagName()
		switch(t) {
		case html.ErrorToken:
			fmt.Println(tz.Err())
			return ""
		case html.TextToken:
			if inbody {
				return string(tz.Text())
			}
		case html.StartTagToken:
			inbody = inbody || string(tn) == "title"
		case html.EndTagToken:
			inbody = inbody && string(tn) != "title"
		}
	}

	return ""
}

// Used to remove the url from the eventual title to avoid botloops
func stripUrls(title string) string {
	return httpRe.ReplaceAllString(title, "")
}

// Pass a url through the is.gd url shortener
func shorten(url string) (string, error) {
	resp, err := client.Get(shortenURL + url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	buf := make([]byte, byteLimit)
	n, _ := io.ReadFull(resp.Body, buf)

	return string(buf[:n]), nil
}
