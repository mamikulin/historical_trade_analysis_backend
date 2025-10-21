package repository

import (
	"archpath/internal/app/models"
	"fmt"

	"gorm.io/gorm"
)

// --- TRADE ANALYSIS (formerly SiteCart) CRUD METHODS ---

// GetAnalysisByUser retrieves the first TradeAnalysis record for a user with a specific status.
func (r *Repository) GetAnalysisByUser(userID uint, status string) (models.TradeAnalysis, error) {
	var analysis models.TradeAnalysis
	err := r.DB.
		Preload("Entries.Artifact").
		Preload("Creator").
		Where("creator_id = ? AND status = ?", userID, status).
		First(&analysis).Error

	if err != nil {
		return models.TradeAnalysis{}, err
	}
	return analysis, nil
}

// GetAnalysisByID retrieves a single TradeAnalysis record by its ID.
func (r *Repository) GetAnalysisByID(analysisID uint) (models.TradeAnalysis, error) {
	var analysis models.TradeAnalysis
	err := r.DB.
		Preload("Entries.Artifact").
		Preload("Creator").
		First(&analysis, analysisID).Error

	if err != nil {
		return models.TradeAnalysis{}, err
	}
	return analysis, nil
}

// CreateAnalysis creates a new TradeAnalysis record.
func (r *Repository) CreateAnalysis(analysis *models.TradeAnalysis) error {
	return r.DB.Create(analysis).Error
}

// UpdateAnalysisDetails updates the name of the archaeological site for a TradeAnalysis record.
func (r *Repository) UpdateAnalysisDetails(analysisID uint, siteName string) error {
	updates := map[string]interface{}{
		"site_name": siteName,
	}

	res := r.DB.Model(&models.TradeAnalysis{}).Where("id = ?", analysisID).Updates(updates)

	if res.Error != nil {
		return fmt.Errorf("error updating analysis details: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// UpdateAnalysisStatusSQL updates the status of a TradeAnalysis record using raw SQL.
func (r *Repository) UpdateAnalysisStatusSQL(analysisID uint, newStatus string) error {
	// Note: table name is typically pluralized snake_case: trade_analyses
	query := `UPDATE trade_analyses SET status = $1 WHERE id = $2`

	result, err := r.SQLDB.Exec(query, newStatus, analysisID)
	if err != nil {
		return fmt.Errorf("SQL UPDATE execution error: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// --- ANALYSIS ARTIFACT RECORD (formerly SiteEntry) CRUD METHODS ---

// GetAnalysisRecord retrieves a single AnalysisArtifactRecord by its composite key.
func (r *Repository) GetAnalysisRecord(analysisID uint, artifactID uint) (models.AnalysisArtifactRecord, error) {
	var record models.AnalysisArtifactRecord
	err := r.DB.
		Where("request_id = ? AND artifact_id = ?", analysisID, artifactID).
		First(&record).Error
	return record, err
}

// CreateAnalysisRecord creates a new AnalysisArtifactRecord.
func (r *Repository) CreateAnalysisRecord(record *models.AnalysisArtifactRecord) error {
	return r.DB.Create(record).Error
}

// UpdateAnalysisRecord updates an existing AnalysisArtifactRecord.
func (r *Repository) UpdateAnalysisRecord(record *models.AnalysisArtifactRecord) error {
	return r.DB.Save(record).Error
}

// UpdateArtifactQuantityInAnalysis updates the quantity of a specific artifact within an analysis record.
func (r *Repository) UpdateArtifactQuantityInAnalysis(analysisID uint, artifactID uint, quantity int) error {
	res := r.DB.Model(&models.AnalysisArtifactRecord{}).
		Where("request_id = ? AND artifact_id = ?", analysisID, artifactID).
		Update("quantity", quantity) // Note: field name changed from artifact_quantity to quantity

	if res.Error != nil {
		return fmt.Errorf("error updating artifact quantity: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// RemoveArtifactFromAnalysis deletes an AnalysisArtifactRecord by its composite key.
func (r *Repository) RemoveArtifactFromAnalysis(analysisID uint, artifactID uint) error {
	res := r.DB.Where("request_id = ? AND artifact_id = ?", analysisID, artifactID).Delete(&models.AnalysisArtifactRecord{})

	if res.Error != nil {
		return fmt.Errorf("error deleting analysis record: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *Repository) UpdateAnalysis(analysis *models.TradeAnalysis) error {
	result := r.DB.Save(analysis)
	if result.Error != nil {
		// Используем logrus для более детального логирования ошибки перед возвратом
		// Если logrus недоступен, можно использовать log.Printf
		return fmt.Errorf("failed to update TradeAnalysis with ID %d: %w", analysis.ID, result.Error)
	}
	return nil
}