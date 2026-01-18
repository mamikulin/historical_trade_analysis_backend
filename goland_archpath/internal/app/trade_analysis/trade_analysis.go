package trade_analysis

import (
	"time"

)

type TradeAnalysis struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	Status         string     `gorm:"not null;default:'draft'" json:"status"`
	CreatorID      uint       `gorm:"not null" json:"creator_id"`
	SiteName       string     `gorm:"not null" json:"site_name"`
	FormationDate  *time.Time `json:"formation_date"`
	CompletionDate *time.Time `json:"completion_date"`
	ModeratorID    *uint      `json:"moderator_id"`
	CompletedEntriesCount int64 `gorm:"default:0" json:"completed_entries_count"`
	CalculatedEntriesCount int64 `gorm:"-" json:"calculated_entries_count"`
}