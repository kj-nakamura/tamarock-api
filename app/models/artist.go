package models

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ArtistInfo is table
type ArtistInfo struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ArtistId  string         `gorm:"not null" json:"artist_id"`
	Name      string         `gorm:"not null" json:"name"`
	Url       string         `json:"url"`
	TwitterId string         `json:"twitter_id"`
	Articles  []Article      `gorm:"many2many:article_artist_infos;" json:"articles"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func migrateArtistInfo() {
	DbConnection.AutoMigrate(&ArtistInfo{})
}

// CreateArtistInfo is アーティスト作成
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

// UpdateArtistInfo is アーティスト更新
func UpdateArtistInfo(r *http.Request, id int) ArtistInfo {
	var artistInfo ArtistInfo

	dec := json.NewDecoder(r.Body)
	for err := dec.Decode(&artistInfo); err != nil && err != io.EOF; {
		log.Println("ERROR: " + err.Error())
	}

	result := DbConnection.Model(&artistInfo).Where("id = ?", id).Updates(ArtistInfo{
		ArtistId:  artistInfo.ArtistId,
		Name:      artistInfo.Name,
		Url:       artistInfo.Url,
		TwitterId: artistInfo.TwitterId,
	})

	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return artistInfo
}

// DeleteArtistInfo is アーティスト情報削除
func DeleteArtistInfo(id int) {
	var artistInfo ArtistInfo

	DbConnection.Delete(&artistInfo, id)
}

// GetArtistInfo is 引数のIDに合致したアーティストを返す
func GetArtistInfo(ID int) ArtistInfo {
	var artistInfo ArtistInfo
	DbConnection.First(&artistInfo, ID)

	return artistInfo
}

// GetArtistInfoFromArtistID is artist_idに合致したアーティストを返す
func GetArtistInfoFromArtistID(artistID string) ArtistInfo {
	var artistInfo ArtistInfo
	var articles []Article

	DbConnection.Where("artist_id = ?", artistID).First(&artistInfo)
	// 記事情報取得
	DbConnection.Debug().Model(&artistInfo).Association("Articles").Find(&articles)

	responseArtistInfo := ArtistInfo{
		ID:        artistInfo.ID,
		ArtistId:  artistInfo.ArtistId,
		Name:      artistInfo.Name,
		Url:       artistInfo.Url,
		TwitterId: artistInfo.TwitterId,
		Articles:  articles,
	}

	return responseArtistInfo
}

// GetArtistInfos is アーティスト情報を複数返す
func GetArtistInfos(start int, end int, order string, sort string, query string) []ArtistInfo {
	var artistInfos []ArtistInfo

	if end > 0 {
		sortColumn := sort
		if sort == "" {
			sortColumn = "id"
		}
		createdOrder := sortColumn + " asc"
		if order == "DESC" {
			createdOrder = sortColumn + " desc"
		}
		limit := end - start
		DbConnection.Order(createdOrder).Offset(start).Limit(limit).Where("name LIKE?", "%"+query+"%").Find(&artistInfos)
	} else {
		DbConnection.Find(&artistInfos)
	}

	return artistInfos
}

// CountArtistInfos is 全記事数を取得
func CountArtistInfos(query string) int {
	var artistInfos []ArtistInfo
	DbConnection.Where("name LIKE?", "%"+query+"%").Find(&artistInfos)

	return len(artistInfos)
}
