# Linkbot
This silly bot just posts links from the mumble chat to a subreddit, that is all.  Does minor special handling for youtube and imgur. 

# TODO
- Fix threading bugs (lots of races)
- Add built in ducking (when bot hears audio lowers volume)
- Extract actual information from songs in list and for playing
- Ability to save playlist and reuse later

# Get Linkbot
    go get -u github.com/robertmeta/linkbot

# Example Use
    linkbot -server="foo.us:64738" -username="foobielinkbot" -password="foo" -insecure=true -reddituser="foouser" -redditpassword="foouserpw" -subreddit="allthefoos"

