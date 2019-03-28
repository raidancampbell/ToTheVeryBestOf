# To The Very Best Of

### A tool to discover the very best tracks an artist has to offer

## Sample Usage:

```
   user@host: go build main.go
   
   user@host: ./main
   
   Enter an artist name:
   barbra streisand
   
   Top tracks for: barbra streisand
   1: Woman in Love, plays: 716433, listeners: 171338
   2: The Way We Were, plays: 288768, listeners: 79421
   3: Memory, plays: 223076, listeners: 62534
   4: Don't Rain on My Parade, plays: 233756, listeners: 58988
   5: Guilty, plays: 150898, listeners: 47704
   6: People, plays: 123115, listeners: 37827
   7: Somewhere, plays: 76138, listeners: 26663
   8: Evergreen, plays: 79255, listeners: 26498
   9: Send in the Clowns, plays: 72998, listeners: 26094
   10: Papa, Can You Hear Me?, plays: 83145, listeners: 24285
   https://www.youtube.com/watch_videos?video_ids=a8DE5U6npkQ,uBPQT2Ia8fU,MWoQW-b6Ph8,-Yfh_CpA9Sk,nVyeNZCENZA
```

## Development Requirements

 - A Last.FM API key, which can be created [here](https://www.last.fm/api/account/create), set with the environment variable `LASTFM_API_KEY`
 - A Youtube Data v3 API key, which can be created [here](https://console.developers.google.com/apis/api/youtube.googleapis.com/overview), set with the environment variable `YOUTUBE_API_KEY`
 
 ## Design Considerations
 
 - Youtube Data API call to make a search costs 100 units, and the free plan is capped at 10,000 units per day.  This gives a maximum of 100 searches per day which is too restrictive at 1 search per track
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
 - Caching / storing youtube search results may ease strain, but may also prove ineffective or too stale
 - Around its peak, ToTheBestOf.com received [over 100,000 requests per week](https://web.archive.org/web/20150301154153/http://tothebestof.com/stats)

 ## Current status
 
 - [X] call Last.FM for track popularity
 - [X] call youtube to use as a song provider
 - [X] basic interactive executable
 - [ ] web JSON API to receive/respond
 - [ ] bundle / reformat code for remote server deployments
 - [ ] webpage to service the JSON API
 - [ ] investigate caching / storing of API calls