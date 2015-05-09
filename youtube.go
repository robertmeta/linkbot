package main

import (
	"html/template"
	"regexp"

	"github.com/layeh/gumble/gumble"
)

var youtubePattern *regexp.Regexp
var youtubeTemplate *template.Template

const youtubeLinkPattern = `https?://(?:www\.)?(?:youtube\.com/watch\?v=|youtu\.be/|youtube\.com/v/|youtube\.com/v/)([[:alnum:]_\-]+)`

func handleYoutubeLink(client *gumble.Client, who, id string) {
	imgURL := "https://i.ytimg.com/vi/" + id + "/hqdefault.jpg"
	linkURL := "https://www.youtube.com/watch?v=" + id
	title := getTitle(linkURL)
	msg := `<b>YouTube</b><br/><center><a href="` + linkURL + `"><img width="250" src="` + imgURL + `"></img><br/>` + title + `</center></a>`
	postLinkToReddit(client, title, "youtube video", who, linkURL)
	sendMsg(client, msg)
}
