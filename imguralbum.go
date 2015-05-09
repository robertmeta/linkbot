package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/layeh/gumble/gumble"
)

var imgurAlbumPattern *regexp.Regexp

const imgurAlbumLinkPattern = `https?://(?:www|i\.)?(?:imgur\.com/)(?:a/|gallery/)([[:alnum:]]+)`

func handleImgurAlbumLink(client *gumble.Client, who, id string) {
	linkURL := "http://imgur.com/gallery/" + id
	title := getTitle(linkURL)
	images := findImages(linkURL)
	msg := `<b>Album</b><br/><center><a href="` + linkURL + `">`
	for _, imgURL := range images {
		msg += `<br/><img width="250" src="` + imgURL + `"></img>`
	}
	msg += `<br/>` + title + `</center></a>`
	postLinkToReddit(client, title, "image album", who, linkURL)
	sendMsg(client, msg)
}

func findImages(url string) []string {
	imgs := []string{}
	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return []string{}
	}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		if strings.Contains(link, "//i.imgur.com") {
			imgLink := `http:` + link[:len(link)-4] + "m" + link[len(link)-4:]
			if len(imgs) < 3 {
				imgs = append(imgs, imgLink)
			}
		}
	})
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("src")
		if strings.Contains(link, "//i.imgur.com") {
			imgLink := `http:` + link[:len(link)-4] + "m" + link[len(link)-4:]
			if len(imgs) < 3 {
				imgs = append(imgs, imgLink)
			}
		}
	})
	return imgs
}
