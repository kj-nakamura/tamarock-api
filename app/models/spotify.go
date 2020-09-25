package models

import (
	"api/config"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

// SpotifyArtistInfo is an artist infomation
type SpotifyArtistInfo struct {
	SpotidyArtistInfo *spotify.FullArtist `json:"spotify_artist_info"`
	YoutubeIds        []string            `json:"youtube_ids"`
}

// NewSpotifyArtist is 引数に準じたartist情報を返す
func NewSpotifyArtist(spotifyArtistInfo *spotify.FullArtist, youtubeIds []string) *SpotifyArtistInfo {
	return &SpotifyArtistInfo{
		spotifyArtistInfo,
		youtubeIds,
	}
}

func GetClient() spotify.Client {
	config := &clientcredentials.Config{
		ClientID:     config.Env.ApiKey,
		ClientSecret: config.Env.ApiSecret,
		TokenURL:     spotify.TokenURL,
	}
	token, err := config.Token(context.Background())
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	client := spotify.Authenticator{}.NewClient(token)

	return client
}

func GetSpotifyArtist(artistID string) *SpotifyArtistInfo {
	client := GetClient()
	result, err := client.GetArtist(spotify.ID(artistID)) // artistID

	if err != nil {
		log.Fatalf("couldn't get artists: %v", err)
	}
	videos := GetMovies(result.Name, "video", 6)

	return NewSpotifyArtist(result, videos)
}

// GetMovies is 引数のIDに合致したアーティストを返す
func GetMovies(query string, searchType string, maxResults int64) []string {
	client := &http.Client{
		Transport: &transport.APIKey{Key: config.Env.YoutubeApiKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	// Make the API call to YouTube.
	call := service.Search.List([]string{"id", "snippet"}).
		Q(query).
		Type(searchType).
		MaxResults(maxResults)
	response, err := call.Do()
	if err != nil {
		fmt.Println(err)
	}

	// Group video, channel, and playlist results in separate lists.
	var videos []string

	// Iterate through each item and add it to the correct list.
	for _, item := range response.Items {
		switch item.Id.Kind {
		case "youtube#video":
			videos = append(videos, item.Id.VideoId)
		}
	}

	return videos
}
