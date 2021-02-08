package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/raidancampbell/ToTheVeryBestOf/data"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

type topTracks struct {
	yt *youtube.Service
	lastfmKey string
	db *gorm.DB
}

func NewTopTracks(youtubeKey, lastFMKey string, db *gorm.DB) *topTracks {
	yt, err := youtube.NewService(context.Background(), option.WithAPIKey(youtubeKey))
	if err != nil {
		panic(err)
	}

	return &topTracks{
		yt:        yt,
		lastfmKey: lastFMKey,
		db:        db,
	}
}


func (t *topTracks) HandleArtistRequest(c *gin.Context) {
	// call the Last.FM api to get the top tracks for the given artist
	lastFMResp, err := getTopTracks(c.Query("Artist"), t.lastfmKey)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		_, _ = c.Writer.WriteString(err.Error())
		return
	}

	var ids []string
	for i := 0; i < 5 && i <= len(lastFMResp.Toptracks.Track); i++ {
		// create the search string for the given track
		track := lastFMResp.Toptracks.Track[i]
		// to hopefully be more correct for vague track titles, we include the artist name in the search
		searchString := fmt.Sprintf("%s %s", track.Artist.Name, track.Name)

		// youtube API call to return search results
		videoID, err := memoizedYoutubeSearch(t.db, t.yt, searchString)
		if err != nil {
			continue
		}
		// get the first result and add it to the top track videos list
		ids = append(ids, videoID)
	}

	// anonymous playlists don't work in embedded players.  looks like we need to create a playlist.
	// update: creating playlists requires way more attention than I'm willing to give this project, plus it was hacky.
	// service accounts won't work, and getting this application verified sounds annoying.
	// SO. plan C: each of the videos is embedded on its own, with the first set to autoplay.
	c.Writer.WriteHeader(http.StatusOK)
	for i, videoID := range ids {
		autoplay := 0
		if i == 0 {
			autoplay = 1
		}
		embed := fmt.Sprintf(`<iframe id="ytplayer" type="text/html" width="640" height="360" src="https://www.youtube.com/embed/%s?autoplay=%d"frameborder="0"></iframe>`, videoID, autoplay)
		_, _ = c.Writer.WriteString(embed)
	}
}

//getTopTracks returns Last.FM's top tracks for the given artist.
// Note the artist name in the response is what the match was for.
func getTopTracks(artistName, lastfmKey string) (*data.LastFMResp, error) {
	artistName = strings.Replace(artistName, " ", "+", -1)
	url := fmt.Sprintf("http://ws.audioscrobbler.com/2.0/?method=artist.gettoptracks&artist=%s&api_key=%s&format=json", artistName, lastfmKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	r := new(data.LastFMResp)
	err = json.NewDecoder(resp.Body).Decode(r)
	return r, err
}


//memoizedYoutubeSearch checks a persisted database to see if the given search query has been searched before
// if it has, the cached response is used.  Otherwise, a youtube search API call is made and cached for future reference
// assumptions:
// the youtube search API call is expensive (in API token allowance), so caching is valuable
// the same search string used multiple times across months/years will produce the same result
func memoizedYoutubeSearch(db *gorm.DB, service *youtube.Service, query string) (string, error) {
	// try to retrieve from the database first
	var id string
	result := data.YoutubeResult{}
	gormRes := db.Model(&data.YoutubeResult{}).Where(&data.YoutubeResult{Query: query}).First(&result)
	if gormRes.RowsAffected != 0 {
		return result.VideoID, nil
	}

	// wasn't found in the database. run the search
	id, err := youtubeSearch(service, query)
	if err != nil {
		return id, err
	}

	// and store the results before returning
	db.Model(&data.YoutubeResult{}).Save(&data.YoutubeResult{
		Query:   query,
		VideoID: id,
	})
	return id, nil
}

// youtubeSearch executes a youtube v3 search list operation (https://developers.google.com/youtube/v3/docs/search/list)
// this costs 100 units.  My account has 10,000 units per day.  So it's limited to 20 artist searches returning the top 5 songs.
// TODO: keep embedded iframe behind a button so that the search operation is only triggered if the user desires
func youtubeSearch(service *youtube.Service, query string) (string, error) {
	call := service.Search.List([]string{"snippet"}).MaxResults(5).Q(query)
	searchListResponse, err := call.Do()
	if err != nil {
		return "", err
	}
	if len(searchListResponse.Items) == 0 {
		return "", fmt.Errorf("no results for query '%s'", query)
	}

	return searchListResponse.Items[0].Id.VideoId, nil
}
