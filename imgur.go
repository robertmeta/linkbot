package main

import (
	"regexp"

	"github.com/layeh/gumble/gumble"
)

var imgurPattern *regexp.Regexp

const imgurLinkPattern = `https?://(?:www|i\.)?(?:imgur\.com/)([[:alnum:]]+)`

func handleImgurLink(client *gumble.Client, who, id string) {
	url := "http://imgur.com/" + id
	title := getTitle(url)
	postLinkToReddit(client, title, who, url)
	message := gumble.TextMessage{
		Channels: []*gumble.Channel{
			client.Self.Channel,
		},
		Message: `<b>POSTED</b><br/><center><a href="` + url + `"><img width="250" src="http://i.imgur.com/` + id + `m.png"></img><br/>` + title + `</center></a>`,
	}
	client.Send(&message)
}
