package user

import (
    "gorm.io/gorm"
)


type User struct {
	gorm.Model
	Login        string `gorm:"uniqueIndex;not null" json:"login"`
	PasswordHash string `json:"-"`
	IsModerator  bool   `gorm:"not null;default:false" json:"is_moderator"`
}