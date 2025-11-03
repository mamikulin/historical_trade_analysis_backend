package trade_analysis

import (
	"time"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

// CreateRequest creates a new trade analysis request
func (r *Repository) CreateRequest(request *TradeAnalysis) error {
	return r.DB.Create(request).Error
}

// GetRequestByID retrieves a request by ID
func (r *Repository) GetRequestByID(id uint) (*TradeAnalysis, error) {
	var request TradeAnalysis
	err := r.DB.First(&request, id).Error
	if err != nil {
		return nil, err
	}
	return &request, nil
}

// GetDraftByCreatorID retrieves the draft request for a specific creator
func (r *Repository) GetDraftByCreatorID(creatorID uint) (*TradeAnalysis, error) {
	var request TradeAnalysis
	err := r.DB.Where("creator_id = ? AND status = ?", creatorID, "draft").First(&request).Error
	if err != nil {
		return nil, err
	}
	return &request, nil
}

// GetAllRequests retrieves all requests with filters (excluding drafts and deleted)
func (r *Repository) GetAllRequests(status string, startDate, endDate *time.Time) ([]TradeAnalysis, error) {
	var requests []TradeAnalysis
	query := r.DB.Where("status != ? AND deleted_at IS NULL", "draft")
	
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	if startDate != nil {
		query = query.Where("formation_date >= ?", startDate)
	}
	
	if endDate != nil {
		query = query.Where("formation_date <= ?", endDate)
	}
	
	err := query.Find(&requests).Error
	return requests, err
}

// UpdateRequest updates a request
func (r *Repository) UpdateRequest(id uint, updates map[string]interface{}) error {
	return r.DB.Model(&TradeAnalysis{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteRequest soft deletes a request (only if it has formation_date)
func (r *Repository) DeleteRequest(id uint) error {
	return r.DB.Delete(&TradeAnalysis{}, id).Error
}

// GetEntriesWithArtifacts retrieves all entries for a request with artifact details
func (r *Repository) GetEntriesWithArtifacts(requestID uint) ([]AnalysisArtifactRecordWithArtifact, error) {
	var entries []AnalysisArtifactRecordWithArtifact
	err := r.DB.Table("analysis_artifact_records").
		Select("analysis_artifact_records.request_id, analysis_artifact_records.artifact_id, analysis_artifact_records.quantity, analysis_artifact_records.order, artifacts.production_center").
		Joins("LEFT JOIN artifacts ON artifacts.id = analysis_artifact_records.artifact_id").
		Where("analysis_artifact_records.request_id = ?", requestID).
		Order("analysis_artifact_records.order ASC").
		Scan(&entries).Error
	return entries, err
}

// CountEntriesByRequestID counts the number of entries for a request
func (r *Repository) CountEntriesByRequestID(requestID uint) (int64, error) {
	var count int64
	err := r.DB.Table("analysis_artifact_records").
		Where("request_id = ?", requestID).
		Count(&count).Error
	return count, err
}