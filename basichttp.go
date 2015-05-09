package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/layeh/gumble/gumble"
)

var basicHTTPPattern *regexp.Regexp
var basicHTTPTemplate *template.Template
var songQueue = []string{}

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
	location := ""
	isSong := false
	if strings.HasSuffix(lowerURL, ".ogg") {
		isSong = true
		location = downloadFromUrl(url)
		os.Rename(location, location+".ogg")
		songQueue = append(songQueue, location+".ogg") // BUG(rm) racey
	}
	if strings.HasSuffix(lowerURL, ".mp3") {
		isSong = true
		location = downloadFromUrl(url)
		os.Rename(location, location+".mp3")
		songQueue = append(songQueue, location+".mp3") // BUG(rm) racey
	}
	if strings.HasSuffix(lowerURL, ".m4a") {
		isSong = true
		location = downloadFromUrl(url)
		os.Rename(location, location+".m4a")
		songQueue = append(songQueue, location+".m4a") // BUG(rm) racey
	}
	if strings.HasSuffix(lowerURL, ".flac") {
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
		if err := stream.Play(streamLoc); err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		go func() {
		TOP:
			stream.Wait()
			log.Println("Cleaning up: ", streamLoc)
			os.Remove(streamLoc)
			if len(songQueue) > 0 {
				streamLoc, songQueue = songQueue[len(songQueue)-1], songQueue[:len(songQueue)-1]
				log.Println("Playing: ", streamLoc)
				if err := stream.Play(streamLoc); err != nil {
					fmt.Printf("%s\n", err)
					return
				}
				goto TOP
			}
		}()

		fmt.Printf("Playing %s\n", streamLoc)
		msg = `<b>Playing Song</b><br/><center><a href="` + url + `">` + url + `</center><br/>Type <b>stop</b> to halt song.</a>`
	}
	postLinkToReddit(client, title, kind, who, url)
	sendMsg(client, msg)
}
