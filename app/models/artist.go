package models

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// ArtistInfo is table
type ArtistInfo struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	ArtistId  string     `gorm:"primary_key: not null" json:"artist_id"`
	Name      string     `json:"name"`
	Url       string     `json:"url"`
	TwitterId string     `json:"twitter_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func migrateArtistInfo() {
	DbConnection.AutoMigrate(&ArtistInfo{})
}

func CreateArtistInfo(r *http.Request) ArtistInfo {
	var artistInfo ArtistInfo
	dec := json.NewDecoder(r.Body)
	for err := dec.Decode(&artistInfo); err != nil && err != io.EOF; {
		log.Println("ERROR: " + err.Error())
	}

	result := DbConnection.Create(&artistInfo)

	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return artistInfo
}

func UpdateArtistInfo(r *http.Request, id int) ArtistInfo {
	var artistInfo ArtistInfo

	dec := json.NewDecoder(r.Body)
	for err := dec.Decode(&artistInfo); err != nil && err != io.EOF; {
		log.Println("ERROR: " + err.Error())
	}

	result := DbConnection.Table("artist_infos").Where("id = ?", id).Update(map[string]interface{}{
		"artist_id":  artistInfo.ArtistId,
		"name":       artistInfo.Name,
		"url":        artistInfo.Url,
		"twitter_id": artistInfo.TwitterId,
	})

	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return artistInfo
}

func DeleteArtistInfo(id int) {
	var artistInfo ArtistInfo

	DbConnection.Delete(&artistInfo, id)
}

// GetArtist is 引数のIDに合致したアーティストを返す
func GetArtistInfo(ID int) ArtistInfo {
	var artistInfo ArtistInfo
	DbConnection.First(&artistInfo, ID)

	return artistInfo
}

// GetArtistFromArtistID is artist_idに合致したアーティストを返す
func GetArtistInfoFromArtistID(artistID string) ArtistInfo {
	var artistInfo ArtistInfo

	DbConnection.Where("artist_id = ?", artistID).First(&artistInfo)

	return artistInfo
}

// GetArtists is アーティスト情報を複数返す
func GetArtistInfos(start int, end int, order string, sort string, query string) []ArtistInfo {
	var artistInfos []ArtistInfo

	sortColumn := sort
	if sort != "" {
		sortColumn = "id"
	}

	createdOrder := sortColumn + " asc"
	if order == "DESC" {
		createdOrder = sortColumn + " desc"
	}
	if end > 0 && start > 0 {
		limit := end - start
		DbConnection.Order(createdOrder).Offset(start).Limit(limit).Where("name LIKE?", "%"+query+"%").Find(&artistInfos)
	}
	DbConnection.Find(&artistInfos)

	return artistInfos
}

// 全記事数を取得
func CountArtistInfos(query string) int {
	var artistInfos []ArtistInfo
	DbConnection.Where("name LIKE?", "%"+query+"%").Find(&artistInfos)

	return len(artistInfos)
}
