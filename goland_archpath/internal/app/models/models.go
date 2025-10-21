package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// --- Custom Type for JSON Storage ---

// AnalysisResult is a custom type to hold the map[string]float64
// and implement GORM's Valuer and Scanner interfaces for JSON storage (Region -> Import Percentage).
type AnalysisResult map[string]float64

// Value implements the driver.Valuer interface, converting the map to a JSON byte array.
func (a AnalysisResult) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(a)
}

// Scan implements the sql.Scanner interface, converting the JSON byte array to a map.
func (a *AnalysisResult) Scan(value interface{}) error {
	if value == nil {
		*a = AnalysisResult{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, a)
}

// --- GORM Models ---

// User (Пользователь)
type User struct {
	gorm.Model
	Login        string `gorm:"uniqueIndex;not null" json:"login"`
	PasswordHash string `json:"-"`
	IsModerator  bool   `gorm:"not null;default:false" json:"is_moderator"`
}

// Artifact (Артефакт/Услуга) - Represents an imported artifact category.
type Artifact struct {
	gorm.Model
	Name        string  `gorm:"not null" json:"name"`
	Description string  `gorm:"type:text" json:"description"`
	IsActive    bool    `gorm:"not null;default:true" json:"is_active"`
	ImageURL    *string `json:"image_url"`

	// FIX: Added default:'Unknown' to prevent the NOT NULL migration error.
	ProductionCenter string `gorm:"not null;default:'Unknown'" json:"production_center"`
	ExampleLocation  *string `json:"example_location"`
}

// TradeAnalysis (Заявка) - Represents a record for trade analysis on a site. (Replaces SiteCart)
type TradeAnalysis struct {
	gorm.Model

	// NotNull Required Fields
	Status    string    `gorm:"not null;default:'draft'" json:"status"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	CreatorID uint      `gorm:"not null" json:"creator_id"`
	SiteName  string    `gorm:"not null" json:"site_name"`

	// Relationships
	Creator User `gorm:"foreignKey:CreatorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"creator"`

	// Dates and Moderator for Workflow
	FormationDate  *time.Time `json:"formation_date"`
	CompletionDate *time.Time `json:"completion_date"`
	ModeratorID    *uint      `json:"moderator_id"`
	Moderator      *User      `gorm:"foreignKey:ModeratorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"moderator"`

	TotalFindsQuantity int `gorm:"not null;default:0" json:"total_finds_quantity"`

	// Stores the map as JSON in a TEXT column.
	AnalysisResult AnalysisResult `gorm:"type:jsonb" json:"analysis_result"` 

	// M:M relationship entries (Replaced SiteEntry)
	Entries []AnalysisArtifactRecord `gorm:"foreignKey:RequestID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"entries"`
}

// AnalysisArtifactRecord (М-М Заявка-Артефакт/Услуга) - Links a trade analysis record to specific artifact categories. (Replaces SiteEntry)
type AnalysisArtifactRecord struct {
	// Composite Unique Key
	RequestID  uint `gorm:"primaryKey" json:"request_id"`
	ArtifactID uint `gorm:"primaryKey" json:"artifact_id"`

	// Additional Fields (Quantity replaces ArtifactQuantity)
	Quantity    int    `gorm:"not null;default:1" json:"quantity"`
	Order       int    `gorm:"not null;default:0" json:"order"`
	IsMainEntry bool   `gorm:"not null;default:false" json:"is_main_entry"`
	Comment     string `gorm:"type:text" json:"comment"`

	// Relationships
	Request  TradeAnalysis `gorm:"foreignKey:RequestID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"request"` 
	Artifact Artifact      `gorm:"foreignKey:ArtifactID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"artifact"`
}

// GetPercentageByRegion calculates the result map.
func (ta *TradeAnalysis) GetPercentageByRegion() AnalysisResult { 
	if len(ta.Entries) == 0 {
		return AnalysisResult{}
	}

	regionCounts := make(map[string]int)
	totalImportQuantity := 0

	for _, entry := range ta.Entries {
		region := entry.Artifact.ProductionCenter 
		quantity := entry.Quantity

		regionCounts[region] += quantity
		totalImportQuantity += quantity
	}

	if totalImportQuantity == 0 {
		return AnalysisResult{}
	}

	regionPercentage := make(AnalysisResult)
	for region, count := range regionCounts {
		percentage := (float64(count) / float64(totalImportQuantity)) * 100
		regionPercentage[region] = percentage
	}

	return regionPercentage
}
