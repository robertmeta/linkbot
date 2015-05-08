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
	lowerURL := strings.ToLower(url)
	if strings.HasSuffix(lowerURL, ".jpg") || strings.HasSuffix(lowerURL, ".jpeg") || strings.HasSuffix(lowerURL, ".png") || strings.HasSuffix(lowerURL, ".gif") {
		kind = "image"
		msg = `<b>Image Posted</b><br/><center><a href="` + url + `"><img width="250" src="` + url + `"></img></center></a>`
	}
	nopost = true
	playSong := false
	location := ""
	if strings.HasSuffix(lowerURL, ".ogg") {
		kind = "ogg"
		location = downloadFromUrl(url)
		streamLoc = location + ".ogg"
		playSong = true
	}
	if strings.HasSuffix(lowerURL, ".mp3") {
		kind = "mp3"
		location = downloadFromUrl(url)
		streamLoc = location + ".mp3"
		playSong = true
	}
	if strings.HasSuffix(lowerURL, ".m4a") {
		kind = "m4a"
		location = downloadFromUrl(url)
		streamLoc = location + ".m4a"
		playSong = true
	}
	if strings.HasSuffix(lowerURL, ".flac") {
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
		stream.Volume = 0.3
		os.Rename(location, streamLoc)
		if err := stream.Play(streamLoc); err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		go func() {
			stream.Wait()
			os.Remove(streamLoc)
		}()

		fmt.Printf("Playing %s\n", streamLoc)
		msg = `<b>Song (vol: 30%)</b><br/><center><a href="` + url + `">` + url + `</center><br/>Type <b>stop</b> to halt song.</a>`
	}
	postLinkToReddit(client, title, kind, who, url)
	sendMsg(client, msg)
}
