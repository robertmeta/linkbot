package main

import (
	"log"
	"strings"
	"time"

	"github.com/jzelinskie/geddit"
	"github.com/layeh/gumble/gumble"
)

var session *geddit.LoginSession

func postLinkToReddit(client *gumble.Client, title, kind, who, link string) {
	if session == nil {
		setupSession(client)
	}

	_, err := session.Me()
	if err != nil { // login then
		setupSession(client)
	}
	if strings.Trim(title, "\n\t ") == "" || strings.Trim(title, "\n\t ") == "Imgur" {
		t := time.Now().Local()
		title = who + " posted a(n) " + kind + " at " + t.Format(time.RubyDate)
	} else {
		title = who + " posted a(n) " + kind + ": " + strings.TrimSpace(title)
	}
	submission := geddit.NewLinkSubmission(subreddit, title, link, false, &geddit.Captcha{})
	submission.Resubmit = false // DIE RESUBMIT
	err = session.Submit(submission)
	if err != nil && strings.Contains(err.Error(), "503") {
		sendMumbleMsg(client, who+` reddit is overloaded at the moment, will retry in a minute.`)
		time.Sleep(1 * time.Minute)
		err = session.Submit(submission)
		if err != nil {
			sendMumbleMsg(client, who+` I give up, stupid broken reddit.`)
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
		sendMumbleMsg(client, `Reddit is overloaded at the moment, I can't even log in, will try again in a minute.`)
		time.Sleep(1 * time.Minute)
		session, err = geddit.NewLoginSession(
			redditUser,
			redditPassword,
			"linkbot v1",
		)
		if err != nil {
			sendMumbleMsg(client, `Reddit is overloaded at the moment, giving up on logging in`)
		}
	}

	sendMumbleMsg(client, `Successfully logged into Reddit.`)
}
