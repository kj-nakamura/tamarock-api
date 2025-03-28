package models

import (
	"api/config"
	"context"
	"fmt"
	"net/http"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

// SpotifyArtistInfo is an artist infomation
type SpotifyArtistInfo struct {
	SpotidyArtistInfo *spotify.FullArtist `json:"spotify_artist_info"`
}

// NewSpotifyArtist is 引数に準じたartist情報を返す
func NewSpotifyArtist(spotifyArtistInfo *spotify.FullArtist) *SpotifyArtistInfo {
	return &SpotifyArtistInfo{
		spotifyArtistInfo,
	}
}

// GetClient is クライアント取得
func GetClient() spotify.Client {
	config := &clientcredentials.Config{
		ClientID:     config.Env.ApiKey,
		ClientSecret: config.Env.ApiSecret,
		TokenURL:     spotify.TokenURL,
	}
	token, err := config.Token(context.Background())
	if err != nil {
		fmt.Sprintf("couldn't get token: %v", err)
	}

	client := spotify.Authenticator{}.NewClient(token)

	return client
}

// GetSpotifyArtist is spotifyからアーティスト情報取得
func GetSpotifyArtist(artistID string) *SpotifyArtistInfo {
	client := GetClient()
	result, err := client.GetArtist(spotify.ID(artistID)) // artistID
	if err != nil {
		fmt.Sprintf("couldn't get artists: %v", err)
		return NewSpotifyArtist(nil)
	}

	return NewSpotifyArtist(result)
}

// GetMovies is 引数のIDに合致したアーティストを返す
func GetMovies(query string, searchType string, maxResults int64) []string {
	client := &http.Client{
		Transport: &transport.APIKey{Key: config.Env.YoutubeApiKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		fmt.Sprintf("Error creating new YouTube client: %v", err)
	}

	var movies []string
	call := service.Search.List([]string{"id", "snippet"}).
		Q(query + " MV").
		Type(searchType).
		MaxResults(maxResults)
	response, err := call.Do()
	if err != nil {
		fmt.Printf("youtube call error:%v\n", err)
		return movies
	}

	for _, item := range response.Items {
		switch item.Id.Kind {
		case "youtube#video":
			movies = append(movies, item.Id.VideoId)
		}
	}

	return movies
}
