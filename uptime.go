package main

import (
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/layeh/gumble/gumble"
)

func handleUptime(e gumble.TextMessageEvent) bool {
	if strings.ToLower(e.Message) == "uptime" {
		sendMumbleMsg(e.Client, "I was last restarted "+humanize.Time(startTime))
		return true
	}
	return false
}
