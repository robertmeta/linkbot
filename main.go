package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumbleutil"
)

var slackChannel string
var slackKey string
var redditUser string
var redditPassword string
var subreddit string
var startTime time.Time

func init() {
	basicHTTPPattern = regexp.MustCompile(basicHTTPLinkPattern)
	imgurPattern = regexp.MustCompile(imgurLinkPattern)
	imgurAlbumPattern = regexp.MustCompile(imgurAlbumLinkPattern)
	youtubePattern = regexp.MustCompile(youtubeLinkPattern)

	flag.StringVar(&redditUser, "reddituser", "", "the reddit user to post as")
	flag.StringVar(&redditPassword, "redditpassword", "", "the reddit user password")
	flag.StringVar(&subreddit, "subreddit", "", "the subreddit to post to")
	flag.StringVar(&slackChannel, "slackchannel", "", "the slack channel to read from and post to")
	flag.StringVar(&slackKey, "slackkey", "", "the slack api key")
}

func main() {
	gul := gumbleutil.Listener{
		Connect:     connectEvent,
		TextMessage: textEvent,
	}
	gumbleutil.Main(extraInit, gul)
}

func connectEvent(e *gumble.ConnectEvent) {
	fmt.Printf("linkbot loaded\n")
}

func extraInit(client *gumble.Client) {
	if redditUser == "" || redditPassword == "" || subreddit == "" || slackKey == "" || slackChannel == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	client.Attach(gumbleutil.AutoBitrate)
	startTime = time.Now()
}

func textEvent(e *gumble.TextMessageEvent) {
	if e.Sender == nil {
		return
	}

	if handleUptime(*e) {
		return
	}

	if handleYoutube(*e) {
		return
	}

	if handleImgurAlbum(*e) {
		return
	}

	if handleImgur(*e) {
		return
	}

	if handleHTTP(*e) {
		return
	}
}
