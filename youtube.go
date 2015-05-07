package main

import (
	"html/template"
	"log"
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
	msg := `<b>YouTube Posted</b><br/><center><a href="` + linkURL + `"><img width="250" src="` + imgURL + `"></img><br/>` + title + `</center></a>`
	log.Println(msg)
	message := gumble.TextMessage{
		Channels: []*gumble.Channel{
			client.Self.Channel,
		},
		Message: msg,
	}
	client.Send(&message)
}
