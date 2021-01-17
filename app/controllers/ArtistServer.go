package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"api/app/models"

	"github.com/gorilla/mux"
	"github.com/zmb3/spotify"
)

// web
func searchArtistHandler(w http.ResponseWriter, r *http.Request) {
	artistName := r.URL.Query().Get("artist_name")
	if artistName == "" {
		APIError(w, "No artist_name param", http.StatusBadRequest)
		return
	}

	client := models.GetClient()
	result, err := client.Search(artistName, spotify.SearchTypeArtist) // artistName
	if err != nil {
		log.Fatalf("couldn't get artists: %v", err)
		return
	}

	// json出力
	responseJSON(w, result)
}

func getArtistHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	ID := vars["id"]

	artist := models.GetSpotifyArtist(ID)

	// json出力
	responseJSON(w, artist)
}

func getArtistInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	artistID := vars["id"]

	if artistID == "" {
		APIError(w, "No artist_id param", http.StatusBadRequest)
		return
	}

	start, _ := strconv.Atoi(r.URL.Query().Get("_start"))
	end, _ := strconv.Atoi(r.URL.Query().Get("_end"))
	order := r.URL.Query().Get("_order")
	sort := r.URL.Query().Get("_sort")
	query := r.URL.Query().Get("q")
	artistInfo := models.GetArtistInfoFromArtistID(artistID, start, end, order, sort, query)

	// json出力
	responseJSON(w, artistInfo)
}

func getArtistInfosHandler(w http.ResponseWriter, r *http.Request) {
	start, _ := strconv.Atoi(r.URL.Query().Get("_start"))
	end, _ := strconv.Atoi(r.URL.Query().Get("_end"))
	order := r.URL.Query().Get("_order")
	sort := r.URL.Query().Get("_sort")
	query := r.URL.Query().Get("q")
	artistInfos := models.GetArtistInfos(start, end, order, sort, query)

	// var artists []*spotify.FullArtist
	// for _, artistInfo := range artistInfos {
	// 	client := models.GetClient()
	// 	artist, err := client.GetArtist(spotify.ID(artistInfo.ArtistId))
	// 	models.NewSpotifyArtist(artistInfo.id, artistInfo.ArtistId, artistInfo.CreatedAt, artistInfo.UpdatedAt, artist.)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	artists = append(artists, artist)
	// }

	// json出力
	responseJSON(w, artistInfos)
}

func getArtistInfosCountHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	artistCount := models.CountArtistInfos(query)

	responseJSON(w, artistCount)
}

// admin
func getAdminArtistsHandler(w http.ResponseWriter, r *http.Request) {
	start, _ := strconv.Atoi(r.URL.Query().Get("_start"))
	end, _ := strconv.Atoi(r.URL.Query().Get("_end"))
	order := r.URL.Query().Get("_order")
	sort := r.URL.Query().Get("_sort")
	query := r.URL.Query().Get("q")
	artistInfos := models.GetArtistInfos(start, end, order, sort, query)
	artistCount := models.CountArtistInfos(query)

	w.Header().Set("X-Total-Count", strconv.Itoa(artistCount))
	responseJSON(w, artistInfos)
}

func getAdminArtistHandler(w http.ResponseWriter, r *http.Request) {
	ID, err := getID(w, r)

	if err != nil {
		fmt.Println(err)
	}

	artistInfo := models.GetArtistInfo(ID)

	responseJSON(w, artistInfo)
}

func createArtistHandler(w http.ResponseWriter, r *http.Request) {
	artistInfo := models.CreateArtistInfo(r)

	responseJSON(w, artistInfo)
}

func updateArtistHandler(w http.ResponseWriter, r *http.Request) {
	ID, err := getID(w, r)

	if err != nil {
		fmt.Println(err)
	}

	artistInfo := models.UpdateArtistInfo(r, ID)

	responseJSON(w, artistInfo)
}

func deleteArtistHandler(w http.ResponseWriter, r *http.Request) {
	ID, err := getID(w, r)

	if err != nil {
		fmt.Println(err)
	}

	models.DeleteArtistInfo(ID)
}
