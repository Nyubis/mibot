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
)

const (
	redirectLimit = 10
	titleLimit    = 200
	timeout       = 2500
	byteLimit     = 65536
)

var httpRe = regexp.MustCompile("https?://[^\\s]*")
var titleRe = regexp.MustCompile("<title>\\s*(?P<want>.*)\\s*</title>")

var client = &http.Client{
	CheckRedirect: func(_ *http.Request, via []*http.Request) error {
		if len(via) >= redirectLimit {
			return errors.New(fmt.Sprintf("Stopped after %d redirects", redirectLimit))
		}
		return nil
	},
	Timeout: time.Duration(timeout)*time.Millisecond,
}

func Handle(msg ircmessage.Message) string {
	url := httpRe.FindString(msg.Content)
	if url == "" {
		return ""
	}

	return ircmessage.PrivMsg(msg.Channel, findTitle(url))
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
		return fmt.Sprintf("%s, %s", contentType, parseSize(resp.ContentLength))
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
