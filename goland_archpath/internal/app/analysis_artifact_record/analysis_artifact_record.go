package analysis_artifact_record

import (
	"time"
)

type AnalysisArtifactRecord struct {
	RequestID  uint `gorm:"primaryKey" json:"request_id"`
	ArtifactID uint `gorm:"primaryKey" json:"artifact_id"`

	Quantity    int    `gorm:"not null;default:1" json:"quantity"`
	Order       int    `gorm:"not null;default:0" json:"order"`
	IsMainEntry bool   `gorm:"not null;default:false" json:"is_main_entry"`
	Comment     string `gorm:"type:text" json:"comment"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (AnalysisArtifactRecord) TableName() string {
	return "analysis_artifact_records"
}