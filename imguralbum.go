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

func handleImgurAlbum(e gumble.TextMessageEvent) bool {
	isAlbum := false

	re := regexp.MustCompile(`imgur.com/r/.+?/`)
	if re.MatchString(e.Message) {
		e.Message = re.ReplaceAllString(e.Message, "imgur.com/")
	}

	if strings.Contains(e.Message, "/a/") {
		isAlbum = true
	}

	imgurAlbumMatches := imgurAlbumPattern.FindStringSubmatch(e.Message)
	if len(imgurAlbumMatches) == 2 {
		go handleImgurAlbumLink(e.Client, e.Sender.Name, imgurAlbumMatches[1], isAlbum)
		return true
	}

	return false
}

func handleImgurAlbumLink(client *gumble.Client, who, id string, isAlbum bool) {
	linkURL := "http://imgur.com/gallery/" + id
	if isAlbum {
		linkURL = "http://imgur.com/a/" + id
	}
	secondaryURL := "http://imgur.com/" + id + ".png"
	title := getTitle(linkURL)
	images := findImages(linkURL)
	msg := `<b>Album</b><br/><center><a href="` + linkURL + `">`
	lastURL := ""
	if len(images) > 0 {
		for _, imgURL := range images {
			msg += `<br/><img width="250" src="` + imgURL + `"></img>`
		}
	} else {
		msg += `<br/><img width="250" src="` + secondaryURL + `"></img>`
	}
	msg += `<br/>` + title + `</center></a>`
	sendSlackMsg(who + " posted gallery " + linkURL)
	postLinkToReddit(client, title, "image album", who, linkURL)
	sendMumbleMsg(client, msg)
	lastURL = ""
	if len(images) > 0 {
		for _, imgURL := range images {
			if lastURL != imgURL {
				sendSlackMsg(imgURL)
				lastURL = imgURL
			}
		}
	} else {
		sendSlackMsg(secondaryURL)
	}
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
	return removeDuplicatesUnordered(imgs)
}

func removeDuplicatesUnordered(elements []string) []string {
	encountered := map[string]bool{}

	// Create a map of all unique elements.
	for v := range elements {
		encountered[elements[v]] = true
	}

	// Place all keys from the map into a slice.
	result := []string{}
	for key, _ := range encountered {
		result = append(result, key)
	}
	return result
}
