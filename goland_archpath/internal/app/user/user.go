package user

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Login        string    `gorm:"unique;not null" json:"login"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Role         string    `gorm:"not null;default:'user'" json:"role"` 
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}