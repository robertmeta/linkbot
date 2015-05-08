package main

import (
	"fmt"
	"html/template"
	"os"
	"regexp"
	"strings"

	"github.com/layeh/gumble/gumble"
)

var basicHTTPPattern *regexp.Regexp
var basicHTTPTemplate *template.Template

const basicHTTPLinkPattern = `href="(https?://.*?)"`

func handlebasicHTTPInfo(client *gumble.Client, who, url string) {
	title := getTitle(url)
	kind := "link"
	msg := `<b>Link Posted</b><br/><center><a href="` + url + `">"` + title + `"</a></center>`
	if strings.HasSuffix(url, ".jpg") || strings.HasSuffix(url, ".jpeg") || strings.HasSuffix(url, ".png") || strings.HasSuffix(url, ".gif") {
		kind = "image"
		msg = `<b>Image Posted</b><br/><center><a href="` + url + `"><img width="250" src="` + url + `"></img></center></a>`
	}
	playSong := false
	location := ""
	if strings.HasSuffix(url, ".ogg") {
		kind = "ogg"
		location = downloadFromUrl(url)
		streamLoc = location + ".ogg"
		playSong = true
	}
	if strings.HasSuffix(url, ".mp3") {
		kind = "mp3"
		location = downloadFromUrl(url)
		streamLoc = location + ".mp3"
		playSong = true
	}
	if strings.HasSuffix(url, ".flac") {
		kind = "flac"
		location = downloadFromUrl(url)
		streamLoc = location + ".flac"
		playSong = true
	}
	if playSong {
		if stream.IsPlaying() {
			stream.Stop()
			os.Remove(streamLoc)
		}
		os.Rename(location, streamLoc)
		if err := stream.Play(streamLoc); err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		fmt.Printf("Playing %s\n", streamLoc)
		msg = `<b>Playing Song</b><br/><center><a href="` + url + `">` + url + `</center><br/>Type <b>stop</b> to halt song.</a>`
	}
	postLinkToReddit(client, title, kind, who, url)
	sendMsg(client, msg)
}
