package models

import "time"

// User is table
type User struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
