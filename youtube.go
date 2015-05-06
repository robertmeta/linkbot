package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"time"

	"github.com/layeh/gumble/gumble"
)

var youtubePattern *regexp.Regexp
var youtubeTemplate *template.Template

const youtubeLinkPattern = `https?://(?:www\.)?(?:youtube\.com/watch\?v=|youtu\.be/|youtube\.com/v/|youtube\.com/v/)([[:alnum:]_\-]+)`
const youtubeResponseTemplate = `
<table>
    <tr>
        <td valign="middle">
		<b>POSTED</b>
        </td>
        <td align="center" valign="middle">
            <a href="http://youtu.be/{{ .Data.Id }}">{{ .Data.Title }} ({{ .Data.Duration }})</a>
        </td>
    </tr>
    <tr>
        <td></td>
        <td align="center">
            <a href="http://youtu.be/{{ .Data.Id }}"><img src="{{ .Data.Thumbnail.HqDefault }}" width="250" /></a>
        </td>
    </tr>
</table>`

func handleYoutubeLink(client *gumble.Client, who, id string) {
	type videoInfo struct {
		Data struct {
			Id        string
			Title     string
			Duration  time.Duration
			Thumbnail struct {
				HqDefault string
			}
		}
	}
	var info videoInfo

	// Fetch + parse video info
	url := fmt.Sprintf("http://gdata.youtube.com/feeds/api/videos/%s?v=2&alt=jsonc", id)
	if resp, err := http.Get(url); err != nil {
		return
	} else {
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&info); err != nil {
			return
		}
		info.Data.Duration *= time.Second
		resp.Body.Close()
	}

	postLinkToReddit(info.Data.Title, who, "https://www.youtube.com/watch?v="+info.Data.Id)

	// Create response string
	var buffer bytes.Buffer
	if err := youtubeTemplate.Execute(&buffer, info); err != nil {
		return
	}
	message := gumble.TextMessage{
		Channels: []*gumble.Channel{
			client.Self.Channel,
		},
		Message: buffer.String(),
	}
	client.Send(&message)
}
