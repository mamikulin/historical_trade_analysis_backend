package trade_analysis

import (
	"archpath/internal/app/analysis_artifact_record"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	DB                        *gorm.DB
	analysisArtifactRecordRepo *analysis_artifact_record.Repository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB:                        db,
		analysisArtifactRecordRepo: analysis_artifact_record.NewRepository(db),
	}
}

func (r *Repository) CreateRequest(request *TradeAnalysis) error {
	return r.DB.Create(request).Error
}

func (r *Repository) GetRequestByID(id uint) (*TradeAnalysis, error) {
	var request TradeAnalysis
	err := r.DB.First(&request, id).Error
	if err != nil {
		return nil, err
	}
	return &request, nil
}

func (r *Repository) GetDraftByCreatorID(creatorID uint) (*TradeAnalysis, error) {
	var request TradeAnalysis
	err := r.DB.Where("creator_id = ? AND status = ?", creatorID, "draft").First(&request).Error
	if err != nil {
		return nil, err
	}
	return &request, nil
}

func (r *Repository) GetAllRequests(status string, startDate, endDate *time.Time, creatorID *uint) ([]TradeAnalysis, error) {
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

	if creatorID != nil {
		query = query.Where("creator_id = ?", *creatorID)
	}

	err := query.Find(&requests).Error
	return requests, err
}

func (r *Repository) UpdateRequest(id uint, updates map[string]interface{}) error {
	return r.DB.Model(&TradeAnalysis{}).Where("id = ?", id).Updates(updates).Error
}

func (r *Repository) DeleteRequest(id uint) error {
	return r.DB.Delete(&TradeAnalysis{}, id).Error
}

func (r *Repository) GetEntriesWithArtifacts(requestID uint) ([]analysis_artifact_record.AnalysisArtifactRecord, error) {
	return r.analysisArtifactRecordRepo.GetRecordsByRequestIDWithArtifacts(requestID)
}

func (r *Repository) CountEntriesByRequestID(requestID uint) (int64, error) {
	var count int64
	err := r.DB.Table("analysis_artifact_records").
		Where("request_id = ?", requestID).
		Count(&count).Error
	return count, err
}

func (r *Repository) UpdateAnalysisArtifactRecord(requestID, artifactID uint, updates map[string]interface{}) error {
	return r.analysisArtifactRecordRepo.UpdateRecord(requestID, artifactID, updates)
}