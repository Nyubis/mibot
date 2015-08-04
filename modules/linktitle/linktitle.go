package linktitle

import (
	"regexp"
	"net/http"
	"io"
	"html"
	"fmt"
	"strings"
	"errors"
	"time"

	"github.com/nyubis/mibot/ircmessage"
	"github.com/nyubis/mibot/core"
)

const (
	redirectLimit = 10
	titleLimit    = 200
	timeout       = 2500
	byteLimit     = 65536

	shortenURL    = "http://is.gd/create.php?format=simple&url="
)

var httpRe = regexp.MustCompile("https?://[^\\s]*")
var titleRe = regexp.MustCompile("<title>\\s*(?P<want>.*)\\s*</title>")

var lastURL string

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
}

func Handle(msg ircmessage.Message) {
	url := httpRe.FindString(msg.Content)
	if url == "" {
		return
	}

	lastURL = url
	bot.SendMessage(msg.Channel, "[URL] " + findTitle(url))
}

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

func findTitle(url string) string {
	resp, err := client.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	buf := make([]byte, byteLimit)
	n, _ := io.ReadFull(resp.Body, buf)
	
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		if resp.ContentLength >= 0 {
			return fmt.Sprintf("%s, %s", contentType, parseSize(resp.ContentLength))
		} else {
			return fmt.Sprintf("%s, unknown size", contentType)
		}
	}

	matches := titleRe.FindStringSubmatch(string(buf[:n]))
	if len(matches) < 2 {
		return ""
	}

	title := html.UnescapeString(matches[1])
	return truncate(title)
}

func parseSize(bytes int64) string {
	units := [...]string{ "B", "KiB", "MiB", "GiB", "TiB", "PiB" }
	size := float64(bytes)
	var i int

	for i = 0; size > 1024 && i < len(units)-1; i++ {
		size /= 1024
	}

	return fmt.Sprintf("%.2f %s", size, units[i])
}

func truncate(s string) string {
	if len(s) <= titleLimit {
		return s
	}
	return s[:titleLimit] + "..."
}

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
