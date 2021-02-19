package models

import (
	"fmt"
	"time"
)

type Youtube struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	MovieID   string    `gorm:"not null" json:"movie_id"`
	ArtistID  int       `gorm:"not null" json:"artist_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func migrateYoutube() {
	DbConnection.AutoMigrate(&Youtube{})
}

// youtubeIDをアーティストを紐付けて保存
func createYoutube(movieID string, artistID int) Youtube {
	youtube := Youtube{
		MovieID:  movieID,
		ArtistID: artistID,
	}

	result := DbConnection.Create(&youtube)

	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return youtube
}

// deleteYoutube is Youtube削除
func deleteYoutube(artistID int) {
	DbConnection.Where("artist_id = ?", artistID).Delete(&Youtube{})
}
