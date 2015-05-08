package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumble_ffmpeg"
	"github.com/layeh/gumble/gumbleutil"
)

var redditUser string
var redditPassword string
var subreddit string
var streamLoc string
var stream *gumble_ffmpeg.Stream
var nopost bool

func init() {
	basicHTTPPattern = regexp.MustCompile(basicHTTPLinkPattern)
	imgurPattern = regexp.MustCompile(imgurLinkPattern)
	imgurAlbumPattern = regexp.MustCompile(imgurAlbumLinkPattern)
	youtubePattern = regexp.MustCompile(youtubeLinkPattern)

	flag.StringVar(&redditUser, "reddituser", "", "the reddit user to post as")
	flag.StringVar(&redditPassword, "redditpassword", "", "the reddit user password")
	flag.StringVar(&subreddit, "subreddit", "", "the subreddit to post to")
}

func extraInit(client *gumble.Client) {
	var err error
	if redditUser == "" || redditPassword == "" || subreddit == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	stream, err = gumble_ffmpeg.New(client)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	client.Attach(gumbleutil.AutoBitrate)
}

func connectEvent(e *gumble.ConnectEvent) {
	fmt.Printf("linkbot loaded\n")
}

func textEvent(e *gumble.TextMessageEvent) {
	nopost = false // face condition, accept for now
	if e.Sender == nil {
		return
	}
	if strings.ToLower(e.Message) == "stop" {
		if stream.IsPlaying() {
			stream.Stop()
			os.Remove(streamLoc)
		}
		return
	}
	if strings.Contains(strings.ToLower(e.Message), "no post") {
		nopost = true
	}
	youtubeMatches := youtubePattern.FindStringSubmatch(e.Message)
	if len(youtubeMatches) == 2 {
		go handleYoutubeLink(e.Client, e.Sender.Name, youtubeMatches[1])
		return
	}

	imgurAlbumMatches := imgurAlbumPattern.FindStringSubmatch(e.Message)
	if len(imgurAlbumMatches) == 2 {
		log.Println(`"` + imgurAlbumMatches[1] + `"`)
		go handleImgurAlbumLink(e.Client, e.Sender.Name, imgurAlbumMatches[1])
		return
	}

	imgurMatches := imgurPattern.FindStringSubmatch(e.Message)
	if len(imgurMatches) == 2 {
		log.Println(`"` + imgurMatches[1] + `"`)
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
