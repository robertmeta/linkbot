package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/layeh/gumble/gumble"
)

var imgurPattern *regexp.Regexp

const imgurLinkPattern = `https?://(?:www|i\.)?(?:imgur\.com/)([[:alnum:]]+)`

func handleImgurLink(client *gumble.Client, who, id string) {
	linkURL := "http://imgur.com/" + id
	imgURL := "http://i.imgur.com/" + id + "m.png"
	title := getTitle(linkURL)
	postLinkToReddit(client, title, who, linkURL)
	msg := `<b>Imgur Posted</b><br/><center><a href="` + linkURL + `"><img width="250" src="` + imgURL + `"></img><br/>` + title + `</center></a>`
	log.Println(msg)
	message := gumble.TextMessage{
		Channels: []*gumble.Channel{
			client.Self.Channel,
		},
		Message: msg,
	}
	client.Send(&message)
}

func findFirstImage(url string) string {
	re := regexp.MustCompile(`href="(https?://i\.imgur\.com/.+)"`)

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return ""
	}
	defer response.Body.Close()
	bs, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error while reading", url, "-", err)
		return ""
	}
	imgURL := re.FindString(string(bs))
	imgURL = strings.Replace(imgURL, `href=`, "", 1)
	imgURL = strings.Replace(imgURL, `"`, "", 2)
	imgURL = strings.Replace(imgURL, `\.`, "m.", 1)
	log.Println(imgURL)
	return imgURL
}
