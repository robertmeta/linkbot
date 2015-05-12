package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
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
var volume float32
var startTime time.Time

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

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
	startTime = time.Now()
}

func connectEvent(e *gumble.ConnectEvent) {
	fmt.Printf("linkbot loaded\n")
}

func textEvent(e *gumble.TextMessageEvent) {
	nopost = false // race condition, accept for now
	if e.Sender == nil {
		return
	}
	if strings.Contains(strings.ToLower(e.Message), "no post") {
		nopost = true
	}
	if strings.ToLower(e.Message) == "uptime" {
		sendMsg(e.Client, "I was last restarted "+humanize.Time(startTime))
	}
	if strings.Contains(strings.ToLower(e.Message), "i love you") {
		sendMsg(e.Client, "I love you too, "+e.Sender.Name+"!")
	}
	if strings.ToLower(e.Message) == "stop" {
		if stream.IsPlaying() {
			songQueue = []string{}
			stream.Stop()
		}
		return
	}
	if strings.ToLower(e.Message) == "mute" {
		if stream.IsPlaying() {
			stream.Volume = 0.0
			sendMsg(e.Client, "Volume set to "+strconv.Itoa(int(stream.Volume*100))+"%")
		}
	}
	if strings.ToLower(e.Message) == "list" {
		msg := "<b>Upcoming Songs:</b><br/><ol>"
		for _, v := range songQueue {
			msg += "<li>" + v + "</li>"
		}
		msg += "</ol>"
		sendMsg(e.Client, msg)
	}
	if strings.ToLower(e.Message) == "full" {
		if stream.IsPlaying() {
			stream.Volume = 1.0
			sendMsg(e.Client, "Volume set to "+strconv.Itoa(int(stream.Volume*100))+"%")
		}
	}
	if strings.ToLower(e.Message) == "next" {
		if stream.IsPlaying() {
			stream.Stop()
		}
	}
	msgParts := strings.Split(e.Message, " ")
	for _, v := range msgParts {
		if strings.ToLower(strings.TrimSpace(v)) == "up" {
			if stream.Volume < 2 {
				stream.Volume += .1
			}
			sendMsg(e.Client, "Volume set to "+strconv.Itoa(int(stream.Volume*100))+"%")
		}
		if strings.ToLower(strings.TrimSpace(v)) == "down" {
			if stream.Volume > 0 {
				stream.Volume -= .1
			}
			sendMsg(e.Client, "Volume set to "+strconv.Itoa(int(stream.Volume*100))+"%")
		}
	}

	youtubeMatches := youtubePattern.FindStringSubmatch(e.Message)
	if len(youtubeMatches) == 2 {
		go handleYoutubeLink(e.Client, e.Sender.Name, youtubeMatches[1])
		return
	}

	imgurAlbumMatches := imgurAlbumPattern.FindStringSubmatch(e.Message)
	if len(imgurAlbumMatches) == 2 {
		go handleImgurAlbumLink(e.Client, e.Sender.Name, imgurAlbumMatches[1])
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
