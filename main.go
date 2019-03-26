package main

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"net/http"
	"os"
	"strconv"
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

func main() {
	// initialize the youtube API
	service, err := youtube.NewService(context.Background(), option.WithAPIKey(youtubeKey))

	// get the user input
	var artist string
	fmt.Println("Enter an artist name:")
	fmt.Scanln(&artist)

	// call the Last.FM api to get the top tracks for the given artist
	topTracks := getTopTracks(artist)

	if err == nil {
		// print out the top tracks
		fmt.Printf("Top tracks for: %s\n", artist)
		for trackNumber, track := range topTracks.Toptracks.Track {
			fmt.Printf("%d: %s, plays: %s, listeners: %s\n", trackNumber+1, track.Name, track.Playcount, track.Listeners)
		}

		// get user input for which track they want a link for
		var input string
		fmt.Println("Enter a song number to fetch the video ID for:")
		_, _ = fmt.Scanln(&input)
		i, _ := strconv.Atoi(input)
		i-- // list was 1-indexed for humans, so we must correct

		// create the search string for the given track
		track := topTracks.Toptracks.Track[i]
		// to hopefully be more correct for vague track titles, we include the artist name in the search
		searchString := fmt.Sprintf("%s %s", track.Artist.Name, track.Name)

		// execute the API call to search. we only really care about the first result
		searchListResponse, err := searchListByKeyword(service, "snippet", 5, searchString, "")

		if err != nil {
			panic(err)
		}
		// format & print the response
		fmt.Println("https://www.youtube.com/watch?v=%s", grabFirstResultID(searchListResponse))

	} else {
		panic(err)
	}
}

func getTopTracks(artistName string) *TopTracksResponse {
	url := buildURL(artistName)
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

func buildURL(artistName string) string {
	replaced := strings.Replace(artistName, " ", "+", -1)
	hostname := "ws.audioscrobbler.com"
	function := "artist.gettoptracks"
	return fmt.Sprintf("http://%s/2.0/?method=%s&artist=%s&api_key=%s&format=json", hostname, function, replaced, lastfmKey)
}

func grabFirstResultID(response *youtube.SearchListResponse) string {
	return response.Items[0].Id.VideoId
}

func searchListByKeyword(service *youtube.Service, part string, maxResults int64, q string, typeArgument string) (*youtube.SearchListResponse, error) {
	call := service.Search.List(part)
	if maxResults != 0 {
		call = call.MaxResults(maxResults)
	}
	if q != "" {
		call = call.Q(q)
	}
	if typeArgument != "" {
		call = call.Type(typeArgument)
	}
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
