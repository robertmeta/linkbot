package main

import (
	"flag"
	"regexp"
	"time"

	"github.com/layeh/gumble/gumble_ffmpeg"
	"github.com/layeh/gumble/gumbleutil"
)

var slackChannel string
var slackKey string
var redditUser string
var redditPassword string
var subreddit string
var streamLoc string
var stream *gumble_ffmpeg.Stream
var nopost bool
var volume float32
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
	flag.StringVar(&slackKey, "slackKey", "", "the slack api key")
}

func main() {
	sendSlackMsg()
	gul := gumbleutil.Listener{
		Connect:     connectEvent,
		TextMessage: textEvent,
	}
	gumbleutil.Main(grumbleExtraInit, gul)
}
