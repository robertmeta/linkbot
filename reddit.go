package main

import (
	"log"
	"strings"
	"time"

	"github.com/jzelinskie/geddit"
	"github.com/layeh/gumble/gumble"
)

var session *geddit.LoginSession

func postLinkToReddit(client *gumble.Client, title, who, link string) {
	if session == nil {
		setupSession(client)
	}

	_, err := session.Me()
	if err != nil { // login then
		setupSession(client)
	}
	if strings.Trim(title, "\n\t ") == "" || strings.Trim(title, "\n\t ") == "Imgur" {
		t := time.Now().Local()
		title = who + " posted at " + t.Format(time.RubyDate)
	} else {
		title = who + " posted: " + strings.TrimSpace(title)
	}
	submission := geddit.NewLinkSubmission(subreddit, title, link, false, &geddit.Captcha{})
	err = session.Submit(submission)
	if err != nil && strings.Contains(err.Error(), "503") {
		message := gumble.TextMessage{
			Channels: []*gumble.Channel{
				client.Self.Channel,
			},
			Message: who + ` reddit is overloaded at the moment, will retry in a minute.`,
		}
		client.Send(&message)
		time.Sleep(1 * time.Minute)
		err = session.Submit(submission)
		if err != nil {
			message := gumble.TextMessage{
				Channels: []*gumble.Channel{
					client.Self.Channel,
				},
				Message: who + ` I give up, stupid broken reddit.`,
			}
			client.Send(&message)
			log.Println("FAILED TO POST: ", submission)
			return
		}
	}
	log.Println("POSTED: ", submission)
}

func setupSession(client *gumble.Client) {
	var err error
	session, err = geddit.NewLoginSession(
		redditUser,
		redditPassword,
		"linkbot v1",
	)
	if err != nil {
		message := gumble.TextMessage{
			Channels: []*gumble.Channel{
				client.Self.Channel,
			},
			Message: `Reddit is overloaded at the moment, I can't even log in, will try again in a minute.`,
		}
		client.Send(&message)
		time.Sleep(1 * time.Minute)
		session, err = geddit.NewLoginSession(
			redditUser,
			redditPassword,
			"linkbot v1",
		)
		if err != nil {
			message := gumble.TextMessage{
				Channels: []*gumble.Channel{
					client.Self.Channel,
				},
				Message: `Reddit is overloaded at the moment, giving up on logging in`,
			}
			client.Send(&message)
		}
	}

	message := gumble.TextMessage{
		Channels: []*gumble.Channel{
			client.Self.Channel,
		},
		Message: `Successfully logged into Reddit.`,
	}
	client.Send(&message)
}
