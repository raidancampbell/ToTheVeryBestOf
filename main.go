package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	artist := "the crystal method"
	url := buildURL(artist)
	resp, err := http.Get(url)

	r := new(TopTracksResponse)
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(r)

	if err == nil {
		fmt.Printf("Top tracks for: %s\n", artist)
		for trackNumber, track := range r.Toptracks.Track {
			fmt.Printf("%d: %s, plays: %s, listeners: %s\n", trackNumber + 1, track.Name, track.Playcount, track.Listeners)
		}
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
