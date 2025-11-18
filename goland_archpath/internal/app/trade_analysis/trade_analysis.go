package trade_analysis

import (
	"time"

	"gorm.io/gorm"
)

type TradeAnalysis struct {
	gorm.Model
	Status             string     `gorm:"not null;default:'draft'" json:"status"`
	CreatorID          uint       `gorm:"not null" json:"creator_id"`
	SiteName           string     `gorm:"not null" json:"site_name"`
	FormationDate      *time.Time `json:"formation_date"`
	CompletionDate     *time.Time `json:"completion_date"`
	ModeratorID        *uint      `json:"moderator_id"`
	TotalFindsQuantity int        `gorm:"not null;default:0" json:"total_finds_quantity"`
}