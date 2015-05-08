package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

func downloadFromUrl(url string) string {
	file, err := ioutil.TempFile(os.TempDir(), "mumble")
	defer file.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return ""
	}
	defer response.Body.Close()

	n, err := io.Copy(file, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return ""
	}

	fmt.Println(n, "bytes downloaded.")
	return file.Name()
}

func sendMsg(client *gumble.Client, msg string) {
	log.Println(msg)
	message := gumble.TextMessage{
		Channels: []*gumble.Channel{
			client.Self.Channel,
		},
		Message: msg,
	}
	client.Send(&message)
}
