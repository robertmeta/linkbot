package main

import (
	"flag"
	"fmt"
	"html/template"
	"os"
	"regexp"

	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumbleutil"
)

var redditUser string
var redditPassword string
var subreddit string

func init() {
	var err error
	basicHTTPPattern = regexp.MustCompile(basicHTTPLinkPattern)
	imgurPattern = regexp.MustCompile(imgurLinkPattern)

	youtubePattern = regexp.MustCompile(youtubeLinkPattern)
	youtubeTemplate, err = template.New("root").Parse(youtubeResponseTemplate)
	if err != nil {
		panic(err)
	}
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
	if len(imgurMatches) == 2 {
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
