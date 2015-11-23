package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/layeh/gumble/gumble"
)

func getTitle(url string) string {
	title := ""
	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return ""
	}
	doc.Find("title").Each(func(i int, s *goquery.Selection) {
		title = s.Text()
	})
	return strings.Trim(title, "\n\t ")
}

func sendMumbleMsg(client *gumble.Client, msg string) {
	log.Println(msg)
	message := gumble.TextMessage{
		Channels: []*gumble.Channel{
			client.Self.Channel,
		},
		Message: msg,
	}
	err := client.Send(&message)
	if err != nil {
		fmt.Println("Error while sending message: ", msg, err)
	}
}
