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
	postLinkToReddit(client, title, who, url)
	msg := `<b>Link Posted</b><br/><center><a href="` + url + `">"` + title + `"</a></center>`
	if strings.HasSuffix(url, ".jpg") || strings.HasSuffix(url, ".jpeg") || strings.HasSuffix(url, ".png") || strings.HasSuffix(url, ".gif") {
		msg = `<b>Image Posted</b><br/><center><a href="` + url + `"><img width="250" src="` + url + `"></img></center></a>`
	}
	if strings.HasSuffix(url, ".ogg") {
		if stream.IsPlaying() {
			stream.Stop()
			os.Remove(streamLoc)
		}
		location := downloadFromUrl(url)
		streamLoc = location + ".ogg"
		os.Rename(location, streamLoc)
		if err := stream.Play(streamLoc); err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		fmt.Printf("Playing %s\n", streamLoc)
		msg = `<b>Playing Song</b><br/><center><a href="` + url + `">` + url + `</center><br/>Type <b>stop</b> to halt song.</a>`
	}
	sendMsg(client, msg)
}
