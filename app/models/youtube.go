package models

import "time"

type Youtube struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	MovieId   string    `gorm:"not null" json:"movie_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func migrateYoutube() {
	DbConnection.AutoMigrate(&Youtube{})
}
