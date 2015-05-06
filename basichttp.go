package main

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"

	"github.com/layeh/gumble/gumble"
	"golang.org/x/net/html"
)

var basicHTTPPattern *regexp.Regexp
var basicHTTPTemplate *template.Template

const basicHTTPLinkPattern = `href="(https?://.*?)"`

func handlebasicHTTPInfo(client *gumble.Client, who, url string) {
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
			postLinkToReddit(client, title, who, url)
			message := gumble.TextMessage{
				Channels: []*gumble.Channel{
					client.Self.Channel,
				},
				Message: `<b>POSTED</b><br/><center><a href="` + url + `">"` + title + `"</a></center>`,
			}
			client.Send(&message)
			return
		case html.TextToken:
			if titleDepth > 0 {
				title = string(z.Text())
			}
		case html.StartTagToken, html.EndTagToken:
			tn, _ := z.TagName()
			//log.Println(string(tn))
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
