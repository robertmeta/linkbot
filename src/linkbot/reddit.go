package main

import (
	"log"
	"strings"
	"time"

	"github.com/aggrolite/geddit"
	"github.com/layeh/gumble/gumble"
)

var session *geddit.OAuthSession

func postLinkToReddit(client *gumble.Client, title, kind, who, link string) {
	if session == nil {
		setupSession(client)
	}

	if session == nil {
		log.Fatal("Couldn't build a session")
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
	_, err = session.Submit(submission)
	if err != nil && strings.Contains(err.Error(), "503") {
		sendMumbleMsg(client, who+` reddit is overloaded at the moment, will retry in a minute.`)
		time.Sleep(1 * time.Minute)
		_, err = session.Submit(submission)
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

	session, err = geddit.NewOAuthSession(
		redditClientID,
		redditClientSecret,
		"LinkBot for u/"+redditUser+" v0.1 see source https://github.com/robertmeta/linkbot",
		"",
	)
	if err != nil {
		sendMumbleMsg(client, `Reddit Error: `+err.Error())
		log.Fatal(err)
	}

	// Create new auth token for confidential clients (personal scripts/apps).
	err = session.LoginAuth(redditUser, redditPassword)
	if err != nil {
		sendMumbleMsg(client, `Reddit Error: `+err.Error())
		log.Fatal(err)
	}

	sendMumbleMsg(client, `Successfully logged into Reddit.`)
}
