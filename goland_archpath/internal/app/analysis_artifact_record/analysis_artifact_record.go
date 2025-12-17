package analysis_artifact_record

import (
	"time"
)

type AnalysisArtifactRecord struct {
	RequestID   uint      `gorm:"primaryKey" json:"request_id"`
	ArtifactID  uint      `gorm:"primaryKey" json:"artifact_id"`
	Quantity    int       `gorm:"not null;default:1" json:"quantity"`
	
	// Результат async-расчета (null пока не вычислено)
	CalculatedValue *float64  `gorm:"type:numeric(10,2)" json:"calculated_value,omitempty"`
	
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relation to Artifact
	Artifact Artifact `gorm:"foreignKey:ArtifactID" json:"artifact,omitempty"`
}

func (AnalysisArtifactRecord) TableName() string {
	return "analysis_artifact_records"
}

// Artifact struct (minimal version with only what you need)
type Artifact struct {
	ID               uint   `gorm:"primaryKey" json:"id"`
	ProductionCenter string `json:"production_center"`
}

func (Artifact) TableName() string {
	return "artifacts"
}