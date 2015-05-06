package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/layeh/gumble/gumble"
	"golang.org/x/net/html"
)

var imgurPattern *regexp.Regexp

const imgurLinkPattern = `https?://(?:www|i\.)?(?:imgur\.com/)([[:alnum:]]+)`

func handleImgurLink(client *gumble.Client, who, id string) {
	url := "http://imgur.com/" + id
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()
	var titleDepth int
	var title string
	z := html.NewTokenizer(response.Body)
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			if strings.Trim(title, "\n \t") == "Imgur" {
				title = ""
			}
			postLinkToReddit(title, who, url)
			message := gumble.TextMessage{
				Channels: []*gumble.Channel{
					client.Self.Channel,
				},
				Message: `Posted "` + title + `" to reddit on behalf of ` + who + `.<br><img width="250" src="http://i.imgur.com/` + id + `m.png"></img>`,
			}
			client.Send(&message)
			return
		case html.TextToken:
			if titleDepth > 0 {
				title = string(z.Text())
			}
		case html.StartTagToken, html.EndTagToken:
			tn, _ := z.TagName()
			if string(tn) == "title" {
				if tt == html.StartTagToken {
					titleDepth++
				} else {
					titleDepth--
				}
			}
		}
	}
}
