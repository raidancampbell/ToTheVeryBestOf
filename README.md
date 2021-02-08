# To The Very Best Of

### A tool to discover the best tracks an artist has to offer

This repository contains the source code of [totheverybestof.com](http://totheverybestof.com), a crappy reimplementation of the (now dead) tothebestof.com.

## Development Requirements

 - A Last.FM API key, which can be created [here](https://www.last.fm/api/account/create), set with the environment variable `LASTFM_API_KEY`
 - A Youtube Data v3 API key, which can be created [here](https://console.developers.google.com/apis/api/youtube.googleapis.com/overview), set with the environment variable `YOUTUBE_API_KEY`
 - execute `go run main.go`, and point your browser at `http://localhost:9071`
 
 ## Design Considerations
 
 - Youtube Data API call to make a search costs 100 units, and the free plan is capped at 10,000 units per day.  This gives a maximum of 100 searches per day which is extremely restrictive at 5 searches per artist
 - youtube search results are cached, but this is unlikely to appreciably help.
 - Alternatives for streaming music are bleak:
    - Grooveshark, the backend for the original ToTheBestOf.com, no longer exists
    - Spotify API doesn't allow more than 30 second previews unless the user is logged in
    - Deezer API is limited to 30 second previews
    - Napster API is limited to 30 second previews
    - Google music has no public API
    - Amazon music has no public API
    - Tidal has no public API
    - Soundcloud hasn't accepted new API application requests in years
    - Pandora doesn't offer streaming chosen tracks
 - Around its peak, ToTheBestOf.com received [over 100,000 requests per week](https://web.archive.org/web/20150301154153/http://tothebestof.com/stats)

 ## Current status
 
 - [X] call Last.FM for track popularity
 - [X] call youtube to use as a song provider
 - [X] basic interactive executable
 - [X] bundle / reformat code for remote server deployments
 - [X] webpage for interacting with service
 - [X] caching / storing of youtube API calls
 - [ ] caching / storing of last.FM API calls
 - [ ] actual logging
 - [ ] hide each embedded video behind a button to lower unused youtube search calls