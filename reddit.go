package main

import (
	"log"
	"strings"
	"time"

	"github.com/jzelinskie/geddit"
)

var session *geddit.LoginSession

func postLinkToReddit(title, who, link string) {
	if session == nil {
		setupSession()
	}

	_, err := session.Me()
	if err != nil { // login then
		setupSession()
	}
	if strings.TrimSpace(title) == "" {
		t := time.Now().Local()
		title = who + " posted at " + t.Format(time.RubyDate)
	} else {
		title = who + " posted: " + strings.TrimSpace(title)
	}
	submission := geddit.NewLinkSubmission(subreddit, title, link, false, &geddit.Captcha{})
	err = session.Submit(submission)
	if err != nil {
		log.Fatal("Unable to make a submission: ", err)
	}
	log.Println(submission)
}

func setupSession() {
	var err error
	session, err = geddit.NewLoginSession(
		redditUser,
		redditPassword,
		"gedditAgent v1",
	)
	if err != nil {
		log.Fatal("Unable to make a session", err)
	}
}
