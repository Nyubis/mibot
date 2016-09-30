package linktitle

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/Nyubis/mibot/core"
	"github.com/Nyubis/mibot/ircmessage"
	"github.com/Nyubis/mibot/modules/admin"
	"github.com/Nyubis/mibot/modules/floodcontrol"
	"github.com/Nyubis/mibot/utils"

	"golang.org/x/net/html"
)

const (
	redirectLimit = 10
	titleLimit    = 200
	timeout       = 5000
	byteLimit     = 65536

	shortenURL = "http://is.gd/create.php?format=simple&url="
)

var httpRe = regexp.MustCompile("https?://[^\\s]*")
var domainRe = regexp.MustCompile("https?://([^\\s/]*)")

var lastURL map[string]string
var disabled map[string]bool

var bot *core.Bot
var client = &http.Client{
	CheckRedirect: func(_ *http.Request, via []*http.Request) error {
		if len(via) >= redirectLimit {
			return errors.New(fmt.Sprintf("Stopped after %d redirects", redirectLimit))
		}
		return nil
	},
	Timeout: time.Duration(timeout) * time.Millisecond,
	/* Why do we make a special connection that limits the amount of traffic,
	*  if we already limit how much bytes we read from the response body?
	*  Because some people make HTTP headers that are over a gigabyte long. Thanks, clsr. */
	Transport: &http.Transport{
		Dial: func(network string, addr string) (net.Conn, error) {
			return NewLimitedConn(network, addr, byteLimit*2)
		},
	},
}

func Init(ircbot *core.Bot) {
	bot = ircbot
	lastURL = make(map[string]string)
	disabled = make(map[string]bool)
	LoadCfg()
}

func LoadCfg() {
	for k, v := range core.Config.Channelsettings {
		if k != "global" && v.Disabled {
			disabled[k] = true
		}
	}
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

	lastURL[msg.Channel] = url
	if !disabled[msg.Channel] && !checkblacklist(url, msg.Channel) {
		title := getAndFindTitle(url)
		if title != "" {
			bot.SendMessage(msg.Channel, "[URL] "+title)
		}
	} else {
		fmt.Println("Link detected, but", msg.Channel, "is disabled")
	}
}

// Incoming event for @shorten
func HandleShorten(_ []string, sender string, channel string) {
	if lastURL[channel] == "" {
		fmt.Println("No last url to shorten")
		return
	}
	short, err := shorten(lastURL[channel])
	if err != nil {
		fmt.Println("Failed to shorten url", lastURL[channel])
		fmt.Println(err)
		return
	}
	bot.SendMessage(channel, sender+": "+short)
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
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	req.Header.Set("User-Agent", "mibot_irc_linkreader/1.0")

	resp, err := client.Do(req)
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
		// Urls are stripped to avoid botloops.
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
		switch t {
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

// Check whether a domain is blacklisted in a particular channel
func checkblacklist(url string, channel string) bool {
	domain := domainRe.FindStringSubmatch(url)[1]
	if utils.Contains(core.Config.Channelsettings["global"].Blacklist, domain) {
		return true
	}
	settings, has := core.Config.Channelsettings[channel]
	return has && utils.Contains(settings.Blacklist, domain)
}
