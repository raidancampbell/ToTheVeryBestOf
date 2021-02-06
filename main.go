package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"net/http"
	"os"
	"strings"
)

// requirements:
// YOUTUBE_API_KEY env variable set to a youtube data v3 api key
// LASTFM_API_KEY env variable set to a last.fm API key
// defined up here to clarify the requirements
var (
	youtubeKey = os.Getenv("YOUTUBE_API_KEY")
	lastfmKey  = os.Getenv("LASTFM_API_KEY")
)

var yt *youtube.Service

func init() {
	// initialize the youtube API
	var err error
	yt, err = youtube.NewService(context.Background(), option.WithAPIKey(youtubeKey))
	if err != nil {
		panic(err)
	}
}

func main() {
	r := gin.Default()
	r.GET("/", handleLanding)
	r.GET("/artist", handleArtistRequest)
	r.Run("0.0.0.0:9071")
}

func handleLanding(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write([]byte(landing))
}

func handleArtistRequest(c *gin.Context) {
	artistName := c.Query("Artist")

	// call the Last.FM api to get the top tracks for the given artist
	topTracks := getTopTracks(artistName)

	var ids []string
	// form an anonymous playlist of the top 5 results
	for i := 0; i < 5; i++ {
		// if last.fm provided less than 5 results, handle accordingly
		if i >= len(topTracks.Toptracks.Track) {
			break
		}

		// create the search string for the given track
		track := topTracks.Toptracks.Track[i]
		// to hopefully be more correct for vague track titles, we include the artist name in the search
		searchString := fmt.Sprintf("%s %s", track.Artist.Name, track.Name)

		// youtube API call to return 5 search results
		response, err := searchListByKeyword(yt, "snippet", 5, searchString)
		if err != nil {
			continue
		}
		// get the first result and add it to the top track videos list
		ids = append(ids, response.Items[0].Id.VideoId)
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
		c.Writer.WriteString(embed)
	}
}

func getTopTracks(artistName string) *TopTracksResponse {
	artistName = strings.Replace(artistName, " ", "+", -1)
	url := fmt.Sprintf("http://ws.audioscrobbler.com/2.0/?method=artist.gettoptracks&artist=%s&api_key=%s&format=json", artistName, lastfmKey)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	r := new(TopTracksResponse)
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		panic(err)
	}
	return r
}

// searchListByKeyword executes a youtube v3 search list operation (https://developers.google.com/youtube/v3/docs/search/list)
// this costs 100 units.  My account has 10,000 units per day.  So it's limited to 20 artist searches returning the top 5 songs.
// TODO: keep embedded iframe behind a button so that the search operation is only triggered if the user desires
func searchListByKeyword(service *youtube.Service, part string, maxResults int64, query string) (*youtube.SearchListResponse, error) {
	call := service.Search.List([]string{part})
	if maxResults != 0 {
		call = call.MaxResults(maxResults)
	}
	call = call.Q(query)
	return call.Do()
}

// thanks, https://mholt.github.io/json-to-go/
type TopTracksResponse struct {
	Toptracks struct {
		Track []struct {
			Name       string `json:"name"`
			Playcount  string `json:"playcount"`
			Listeners  string `json:"listeners"`
			Mbid       string `json:"mbid,omitempty"`
			URL        string `json:"url"`
			Streamable string `json:"streamable"`
			Artist     struct {
				Name string `json:"name"`
				Mbid string `json:"mbid"`
				URL  string `json:"url"`
			} `json:"artist"`
			Image []struct {
				Text string `json:"#text"`
				Size string `json:"size"`
			} `json:"image"`
			Attr struct {
				Rank string `json:"rank"`
			} `json:"@attr"`
		} `json:"track"`
		Attr struct {
			Artist     string `json:"artist"`
			Page       string `json:"page"`
			PerPage    string `json:"perPage"`
			TotalPages string `json:"totalPages"`
			Total      string `json:"total"`
		} `json:"@attr"`
	} `json:"toptracks"`
}
