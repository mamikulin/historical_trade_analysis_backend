package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Login        string `gorm:"uniqueIndex;not null"`
	PasswordHash string
	IsModerator  bool `gorm:"not null;default:false"`
}

type Artifact struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Period      string
	Region      string
	Description string `gorm:"type:text"`
	ImageURL    *string
	IsActive    bool `gorm:"not null;default:true"`
}

type SiteCart struct {
	gorm.Model

	Status    string `gorm:"not null;default:'draft'"`
	CreatorID uint   `gorm:"not null"`
	SiteName  string `gorm:"not null"`

	Creator User `gorm:"foreignKey:CreatorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	FormationDate  *time.Time
	CompletionDate *time.Time
	ModeratorID    *uint
	Moderator      *User `gorm:"foreignKey:ModeratorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	CalculatedCost float64 `gorm:"not null;default:0"`

	Entries []SiteEntry `gorm:"foreignKey:CartID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

type SiteEntry struct {
	CartID     uint `gorm:"primaryKey"`
	ArtifactID uint `gorm:"primaryKey"`

	ArtifactQuantity int    `gorm:"not null;default:1"`
	Comment          string `gorm:"type:text"`
	IsMain           bool   `gorm:"not null;default:false"`

	SiteCart SiteCart `gorm:"foreignKey:CartID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Artifact Artifact `gorm:"foreignKey:ArtifactID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (cart *SiteCart) GetPercentageByRegion() map[string]float64 {
	if len(cart.Entries) == 0 {
		return make(map[string]float64)
	}

	regionCounts := make(map[string]int)
	totalArtifacts := 0

	for _, entry := range cart.Entries {
		region := entry.Artifact.Region
		quantity := entry.ArtifactQuantity

		regionCounts[region] += quantity
		totalArtifacts += quantity
	}

	regionPercentage := make(map[string]float64)
	for region, count := range regionCounts {
		percentage := (float64(count) / float64(totalArtifacts)) * 100
		regionPercentage[region] = percentage
	}

	return regionPercentage
}