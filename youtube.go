package main

import (
	"regexp"

	"github.com/layeh/gumble/gumble"
)

var youtubePattern *regexp.Regexp

const youtubeLinkPattern = `https?://(?:www\.)?(?:youtube\.com/watch\?v=|youtu\.be/|youtube\.com/v/|youtube\.com/v/)([[:alnum:]_\-]+)`

func handleYoutube(e gumble.TextMessageEvent) bool {
	re2 := regexp.MustCompile(`.*?t=(.+?)"`)
	timeStart := "0"
	times := re2.FindStringSubmatch(e.Message)
	if len(times) > 0 {
		timeStart = times[1]
	}

	youtubeMatches := youtubePattern.FindStringSubmatch(e.Message)
	if len(youtubeMatches) == 2 {
		go handleYoutubeLink(e.Client, e.Sender.Name, youtubeMatches[1], timeStart)
		return true
	}

	return false
}

func handleYoutubeLink(client *gumble.Client, who, id, t string) {
	imgURL := "https://i.ytimg.com/vi/" + id + "/hqdefault.jpg"
	linkURL := "https://www.youtube.com/watch?v=" + id + "&t=" + t
	title := getTitle(linkURL)
	msg := `<b>YouTube</b><br/><center><a href="` + linkURL + `"><img width="250" src="` + imgURL + `"></img><br/>` + title + `</center></a>`
	postLinkToReddit(client, title, "youtube video", who, linkURL)
	sendMumbleMsg(client, msg)
	sendSlackMsg(who + " posted " + linkURL)
}
