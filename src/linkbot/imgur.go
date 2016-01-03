package main

import (
	"regexp"

	"github.com/layeh/gumble/gumble"
)

var imgurPattern *regexp.Regexp

const imgurLinkPattern = `https?://(?:www|i\.)?(?:imgur\.com/)([[:alnum:]]+)`

func handleImgur(e gumble.TextMessageEvent) bool {
	re := regexp.MustCompile(`imgur.com/r/.+?/`)
	if re.MatchString(e.Message) {
		e.Message = re.ReplaceAllString(e.Message, "imgur.com/")
	}

	imgurMatches := imgurPattern.FindStringSubmatch(e.Message)
	if len(imgurMatches) == 2 {
		go handleImgurLink(e.Client, e.Sender.Name, imgurMatches[1])
		return true
	}
	return false
}

func handleImgurLink(client *gumble.Client, who, id string) {
	linkURL := "http://imgur.com/" + id
	imgURL := "http://i.imgur.com/" + id + "m.png"
	title := getTitle(linkURL)
	postLinkToReddit(client, title, "image", who, linkURL)
	msg := `<b>Imgur</b><br/><center><a href="` + linkURL + `"><img width="250" src="` + imgURL + `"></img><br/>` + title + `</center></a>`
	sendMumbleMsg(client, msg)
	sendSlackMsg(who + " posted " + linkURL)
}
