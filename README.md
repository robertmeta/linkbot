# Linkbot
This silly bot just posts links from the mumble chat to a subreddit, that is all.  Does minor special handling for youtube and imgur. 

# TODO
- Song queue for when a song is playing, rather than just stopping it
- Show upcoming playlist
- Ability to save playlist and reuse later

# Get Linkbot
    go get -u github.com/robertmeta/linkbot

# Example Use
    linkbot -server="foo.us:64738" -username="foobielinkbot" -password="foo" -insecure=true -reddituser="foouser" -redditpassword="foouserpw" -subreddit="allthefoos"
