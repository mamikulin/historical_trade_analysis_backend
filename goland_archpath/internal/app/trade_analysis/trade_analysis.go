package trade_analysis

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
	"gorm.io/gorm"
)

type AnalysisResult map[string]float64

func (a *AnalysisResult) Scan(value interface{}) error {
	if value == nil {
		*a = make(AnalysisResult)
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal AnalysisResult value: %v", value)
	}
	
	result := make(AnalysisResult)
	if err := json.Unmarshal(bytes, &result); err != nil {
		return err
	}
	
	*a = result
	return nil
}

func (a AnalysisResult) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(a)
}

type TradeAnalysis struct {
	gorm.Model
	
	Status    string    `gorm:"not null;default:'draft'" json:"status"`
	CreatorID uint      `gorm:"not null" json:"creator_id"`
	SiteName  string    `gorm:"not null" json:"site_name"`
	
	FormationDate  *time.Time `json:"formation_date"`
	CompletionDate *time.Time `json:"completion_date"`
	ModeratorID    *uint      `json:"moderator_id"`
	
	TotalFindsQuantity int `gorm:"not null;default:0" json:"total_finds_quantity"`
	
	AnalysisResult AnalysisResult `gorm:"type:jsonb" json:"analysis_result"`
}

func (ta *TradeAnalysis) GetPercentageByRegion(entries []AnalysisArtifactRecordWithArtifact) AnalysisResult {
	if len(entries) == 0 {
		return AnalysisResult{}
	}
	
	regionCounts := make(map[string]int)
	totalImportQuantity := 0
	
	for _, entry := range entries {
		region := entry.ProductionCenter
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

type AnalysisArtifactRecordWithArtifact struct {
	RequestID        uint   `json:"request_id"`
	ArtifactID       uint   `json:"artifact_id"`
	Quantity         int    `json:"quantity"`
	Order            int    `json:"order"`
	ProductionCenter string `json:"production_center"`
}