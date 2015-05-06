package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumbleutil"
	"golang.org/x/net/html"
)

var redditUser string
var redditPassword string
var subreddit string

func init() {
	basicHTTPPattern = regexp.MustCompile(basicHTTPLinkPattern)
	imgurPattern = regexp.MustCompile(imgurLinkPattern)
	youtubePattern = regexp.MustCompile(youtubeLinkPattern)

	flag.StringVar(&redditUser, "reddituser", "", "the reddit user to post as")
	flag.StringVar(&redditPassword, "redditpassword", "", "the reddit user password")
	flag.StringVar(&subreddit, "subreddit", "", "the subreddit to post to")
}

func extraInit(client *gumble.Client) {
	if redditUser == "" || redditPassword == "" || subreddit == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func connectEvent(e *gumble.ConnectEvent) {
	fmt.Printf("linkbot loaded\n")
}

func textEvent(e *gumble.TextMessageEvent) {
	if e.Sender == nil {
		return
	}
	youtubeMatches := youtubePattern.FindStringSubmatch(e.Message)
	if len(youtubeMatches) == 2 {
		go handleYoutubeLink(e.Client, e.Sender.Name, youtubeMatches[1])
		return
	}

	imgurMatches := imgurPattern.FindStringSubmatch(e.Message)
	if len(imgurMatches) == 2 && imgurMatches[1] != "gallery" {
		go handleImgurLink(e.Client, e.Sender.Name, imgurMatches[1])
		return
	}

	basicHTTPMatches := basicHTTPPattern.FindStringSubmatch(e.Message)
	if len(basicHTTPMatches) == 2 {
		go handlebasicHTTPInfo(e.Client, e.Sender.Name, basicHTTPMatches[1])
		return
	}
}

func main() {
	gul := gumbleutil.Listener{
		Connect:     connectEvent,
		TextMessage: textEvent,
	}
	gumbleutil.Main(extraInit, gul)
}

func getTitle(url string) string {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return ""
	}
	defer response.Body.Close()
	var titleDepth int
	var title string
	z := html.NewTokenizer(response.Body)
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return title
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
	return title
}
