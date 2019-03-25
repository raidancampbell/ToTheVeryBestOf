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

func main() {
	artist := "the crystal method"
	url := buildURL(artist)
	resp, err := http.Get(url)

	r := new(TopTracksResponse)
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(r)

	ctx := context.Background()
	service, _ := youtube.NewService(ctx, option.WithAPIKey(os.Getenv("YOUTUBE_API_KEY")))

	if err == nil {
		fmt.Printf("Top tracks for: %s\n", artist)
		for trackNumber, track := range r.Toptracks.Track {
			fmt.Printf("%d: %s, plays: %s, listeners: %s\n", trackNumber+1, track.Name, track.Playcount, track.Listeners)
		}
		var input string
		fmt.Println("Enter a song number to fetch the video ID for:")
		_, _ = fmt.Scanln(&input)
		i, _ := strconv.Atoi(input)

		searchListResponse, _ := searchListByKeyword(service, "snippet", 25, r.Toptracks.Track[i-1].Name, "")
		fmt.Println(grabFirstResultID(searchListResponse))

	} else {
		panic(err)
	}
}

func buildURL(artistName string) string {
	replaced := strings.Replace(artistName, " ", "+", -1)
	hostname := "ws.audioscrobbler.com"
	function := "artist.gettoptracks"
	return fmt.Sprintf("http://%s/2.0/?method=%s&artist=%s&api_key=%s&format=json", hostname, function, replaced, getAPIKey())
}

func getAPIKey() string {
	return os.Getenv("LASTFM_API_KEY")
}

func pprintJson(inBytes interface{}) {
	formattedBytes, _ := json.MarshalIndent(inBytes, "", "  ")
	os.Stdout.Write(formattedBytes)
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
