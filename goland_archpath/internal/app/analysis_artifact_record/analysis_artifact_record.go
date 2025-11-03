package analysis_artifact_record

import (
	"time"
)

// AnalysisArtifactRecord (М-М Заявка-Артефакт/Услуга) - Links a trade analysis record to specific artifact categories.
type AnalysisArtifactRecord struct {
	// Composite Primary Key
	RequestID  uint `gorm:"primaryKey" json:"request_id"`
	ArtifactID uint `gorm:"primaryKey" json:"artifact_id"`

	// Additional Fields
	Quantity    int    `gorm:"not null;default:1" json:"quantity"`
	Order       int    `gorm:"not null;default:0" json:"order"`
	IsMainEntry bool   `gorm:"not null;default:false" json:"is_main_entry"`
	Comment     string `gorm:"type:text" json:"comment"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for GORM
func (AnalysisArtifactRecord) TableName() string {
	return "analysis_artifact_records"
}