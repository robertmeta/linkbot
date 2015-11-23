package main

import (
	"regexp"
	"strings"

	"github.com/layeh/gumble/gumble"
)

var basicHTTPPattern *regexp.Regexp

const basicHTTPLinkPattern = `href="(https?://.*?)"`

func handleHTTP(e gumble.TextMessageEvent) bool {
	basicHTTPMatches := basicHTTPPattern.FindStringSubmatch(e.Message)
	if len(basicHTTPMatches) == 2 {
		go handlebasicHTTPInfo(e.Client, e.Sender.Name, basicHTTPMatches[1])
		return true
	}
	return false
}

func handlebasicHTTPInfo(client *gumble.Client, who, url string) {
	title := getTitle(url)
	kind := "link"
	msg := `<b>Link</b><br/><center><a href="` + url + `">"` + title + `"</a></center>`
	lowerURL := strings.ToLower(url)
	if strings.HasSuffix(lowerURL, ".jpg") || strings.HasSuffix(lowerURL, ".jpeg") || strings.HasSuffix(lowerURL, ".png") || strings.HasSuffix(lowerURL, ".gif") {
		kind = "image"
		msg = `<b>Image Posted</b><br/><center><a href="` + url + `"><img width="250" src="` + url + `"></img></center></a>`
	}
	postLinkToReddit(client, title, kind, who, url)
	sendMumbleMsg(client, msg)
	sendSlackMsg(who + " posted " + url)
}
