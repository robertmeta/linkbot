package main

import (
	"log"

	"github.com/bluele/slack"
)

func sendSlackMsg(msg string) {
	api := slack.New(slackKey)
	channel, err := api.FindChannelByName(slackChannel)
	if err != nil {
		log.Fatal(err)
	}
	err = api.ChatPostMessage(channel.Id, msg, nil)
	if err != nil {
		log.Fatal(err)
	}
}
