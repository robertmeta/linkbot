package main

import (
	"fmt"
	"html/template"
	"os"
	"regexp"
	"strings"

	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumble_ffmpeg"
)

var basicHTTPPattern *regexp.Regexp
var basicHTTPTemplate *template.Template
var songQueue = []string{}

const basicHTTPLinkPattern = `href="(https?://.*?)"`

func handlebasicHTTPInfo(client *gumble.Client, who, url string) {
	title := getTitle(url)
	kind := "link"
	msg := `<b>Link</b><br/><center><a href="` + url + `">"` + title + `"</a></center>`
	lowerURL := strings.ToLower(url)
	if strings.HasSuffix(lowerURL, ".jpg") || strings.HasSuffix(lowerURL, ".jpeg") || strings.HasSuffix(lowerURL, ".png") || strings.HasSuffix(lowerURL, ".gif") {
		kind = "image"
		msg = `<b>Image Posted</b><br/><center><a href="` + url + `"><img width="250" src="` + url + `"></img></center></a>`
	}
	location := ""
	isSong := false
	if strings.HasSuffix(lowerURL, ".ogg") {
		nopost = true
		isSong = true
		location = downloadFromUrl(url)
		os.Rename(location, location+".ogg")
		songQueue = append(songQueue, location+".ogg") // BUG(rm) racey
	}
	if strings.HasSuffix(lowerURL, ".mp3") {
		nopost = true
		isSong = true
		location = downloadFromUrl(url)
		os.Rename(location, location+".mp3")
		songQueue = append(songQueue, location+".mp3") // BUG(rm) racey
	}
	if strings.HasSuffix(lowerURL, ".m4a") {
		nopost = true
		isSong = true
		location = downloadFromUrl(url)
		os.Rename(location, location+".m4a")
		songQueue = append(songQueue, location+".m4a") // BUG(rm) racey
	}
	if strings.HasSuffix(lowerURL, ".flac") {
		nopost = true
		isSong = true
		location = downloadFromUrl(url)
		os.Rename(location, location+".flac")
		songQueue = append(songQueue, location+".flac") // BUG(rm) racey
	}
	if isSong && len(songQueue) > 0 {
		if stream.IsPlaying() {
			sendMsg(client, `Queued <br><center><a href="`+url+`">`+url+`</a></center>`)
			return
		}
		stream.Volume = 0.2
		streamLoc, songQueue = songQueue[len(songQueue)-1], songQueue[:len(songQueue)-1]
		stream.Source = gumble_ffmpeg.SourceFile(streamLoc)
		if err := stream.Play(); err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		go func() {
			for streamLoc != "" {
				stream.Wait()
				os.Remove(streamLoc)
				if len(songQueue) > 0 {
					streamLoc, songQueue = songQueue[len(songQueue)-1], songQueue[:len(songQueue)-1]
					stream.Source = gumble_ffmpeg.SourceFile(streamLoc)
					if err := stream.Play(); err != nil {
						fmt.Printf("%s\n", err)
						os.Remove(streamLoc)
						return
					}
				} else {
					break
				}
			}
		}()

		fmt.Printf("Playing %s\n", streamLoc)
		msg = `<b>Playing Song</b><br/><center><a href="` + url + `">` + url + `</center><br/>Type <b>stop</b> to halt song.</a>`
	}
	postLinkToReddit(client, title, kind, who, url)
	sendMsg(client, msg)
}
