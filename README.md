# Linkbot
This silly bot just posts links from the mumble chat to a subreddit, and now to a Slack channel. Does minor special handling for youtube and imgur.  Don't expect it to be top
quality, it is my goof around project, written in between league of legends games.

# Notes
The init is wierd because of two sets of flags, those for our app, and those passed through to grumble, so look in grumbleExtraInit for init.

# Get Linkbot
    go get -u github.com/robertmeta/linkbot

# Example Use
    linkbot -server="foo.us:64738" -username="foobielinkbot" -password="foo" -insecure=true -reddituser="foouser" -redditpassword="foouserpw" -subreddit="allthefoos" -slackkey="ASDFASDFASDFASDFAFSDASDFASDF" -slackchannel="general"
