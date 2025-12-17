package analysis_artifact_record

import (
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) CreateRecord(record *AnalysisArtifactRecord) error {
	return r.DB.Create(record).Error
}

func (r *Repository) GetRecordByCompositeKey(requestID, artifactID uint) (*AnalysisArtifactRecord, error) {
	var record AnalysisArtifactRecord
	err := r.DB.Where("request_id = ? AND artifact_id = ?", requestID, artifactID).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *Repository) GetRecordsByRequestID(requestID uint) ([]AnalysisArtifactRecord, error) {
	var records []AnalysisArtifactRecord
	err := r.DB.Where("request_id = ?", requestID).Find(&records).Error
	return records, err
}

func (r *Repository) GetRecordsByRequestIDWithArtifacts(requestID uint) ([]AnalysisArtifactRecord, error) {
	var records []AnalysisArtifactRecord
	err := r.DB.
		Preload("Artifact").
		Where("request_id = ?", requestID).
		Find(&records).Error
	return records, err
}

func (r *Repository) GetRecordsByArtifactID(artifactID uint) ([]AnalysisArtifactRecord, error) {
	var records []AnalysisArtifactRecord
	err := r.DB.Where("artifact_id = ?", artifactID).Find(&records).Error
	return records, err
}

func (r *Repository) UpdateRecord(requestID, artifactID uint, updates map[string]interface{}) error {
	return r.DB.Model(&AnalysisArtifactRecord{}).
		Where("request_id = ? AND artifact_id = ?", requestID, artifactID).
		Updates(updates).Error
}

func (r *Repository) DeleteRecord(requestID, artifactID uint) error {
	return r.DB.Where("request_id = ? AND artifact_id = ?", requestID, artifactID).
		Delete(&AnalysisArtifactRecord{}).Error
}

func (r *Repository) DeleteRecordsByRequestID(requestID uint) error {
	return r.DB.Where("request_id = ?", requestID).Delete(&AnalysisArtifactRecord{}).Error
}