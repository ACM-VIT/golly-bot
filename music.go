package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

// For handling dynamic videoId and playlistId keys
type Info map[string]string

type Response struct {
	Kind          string `json:"kind"`
	Etag          string `json:"etag"`
	NextPageToken string `json:"nextPageToken"`
	RegionCode    string `json:"regionCode"`
	PageInfo      struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
	Items []struct {
		Kind    string `json:"kind"`
		Etag    string `json:"etag"`
		ID      Info   `json:"id"`
		Snippet struct {
			PublishedAt time.Time `json:"publishedAt"`
			ChannelID   string    `json:"channelId"`
			Title       string    `json:"title"`
			Description string    `json:"description"`
			Thumbnails  struct {
				Default struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"default"`
				Medium struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"medium"`
				High struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"high"`
			} `json:"thumbnails"`
			ChannelTitle         string    `json:"channelTitle"`
			LiveBroadcastContent string    `json:"liveBroadcastContent"`
			PublishTime          time.Time `json:"publishTime"`
		} `json:"snippet"`
	} `json:"items"`
}

func searchYT(query string, youtubeAPIKey string) (string, string) {
	flag.Parse()

	// URL encode the query
	query = url.QueryEscape(query)

	// Make the API call to YouTube.
	resp, err := http.Get("https://youtube.googleapis.com/youtube/v3/search?part=snippet&maxResults=1&q=" + query + "&key=" + youtubeAPIKey)
	if err != nil {
		log.Fatalf("Error making search API call: %v", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var result Response
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	items := result.Items
	// Check if query returned any results
	if len(items) == 0 {
		return "No search results were found!", ""
	}

	// Check if result is video or playlist
	kind := items[0].ID["kind"]
	var link string
	if kind == "youtube#video" {
		videoId := items[0].ID["videoId"]
		link = "https://www.youtube.com/watch?v=" + videoId
	} else {
		playlistId := items[0].ID["playlistId"]
		link = "https://youtube.com/playlist?list=" + playlistId
	}

	// Parse title and replace html entities
	title := items[0].Snippet.Title
	title = html.UnescapeString(title)

	return title, link
}
