package models

import "time"

// Category is table
type Category struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	Name      string     `gorm:"not null" json:"name"`
	Articles  []Article  `json:"articles"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func migrateCategory() {
	// DbConnection.AutoMigrate(&Category{})
}
