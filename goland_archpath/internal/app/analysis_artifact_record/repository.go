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

// CreateRecord creates a new analysis-artifact record
func (r *Repository) CreateRecord(record *AnalysisArtifactRecord) error {
	return r.DB.Create(record).Error
}

// GetRecordByCompositeKey retrieves a record by its composite key
func (r *Repository) GetRecordByCompositeKey(requestID, artifactID uint) (*AnalysisArtifactRecord, error) {
	var record AnalysisArtifactRecord
	err := r.DB.Where("request_id = ? AND artifact_id = ?", requestID, artifactID).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// GetRecordsByRequestID retrieves all records for a specific request
func (r *Repository) GetRecordsByRequestID(requestID uint) ([]AnalysisArtifactRecord, error) {
	var records []AnalysisArtifactRecord
	err := r.DB.Where("request_id = ?", requestID).Order("`order` ASC").Find(&records).Error
	return records, err
}

// GetRecordsByArtifactID retrieves all records for a specific artifact
func (r *Repository) GetRecordsByArtifactID(artifactID uint) ([]AnalysisArtifactRecord, error) {
	var records []AnalysisArtifactRecord
	err := r.DB.Where("artifact_id = ?", artifactID).Find(&records).Error
	return records, err
}

// UpdateRecord updates an existing record by composite key
func (r *Repository) UpdateRecord(requestID, artifactID uint, updates map[string]interface{}) error {
	return r.DB.Model(&AnalysisArtifactRecord{}).
		Where("request_id = ? AND artifact_id = ?", requestID, artifactID).
		Updates(updates).Error
}

// DeleteRecord deletes a record by composite key
func (r *Repository) DeleteRecord(requestID, artifactID uint) error {
	return r.DB.Where("request_id = ? AND artifact_id = ?", requestID, artifactID).
		Delete(&AnalysisArtifactRecord{}).Error
}

// DeleteRecordsByRequestID deletes all records for a specific request
func (r *Repository) DeleteRecordsByRequestID(requestID uint) error {
	return r.DB.Where("request_id = ?", requestID).Delete(&AnalysisArtifactRecord{}).Error
}