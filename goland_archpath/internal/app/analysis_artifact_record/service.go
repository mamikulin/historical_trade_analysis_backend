package analysis_artifact_record

import (
	"fmt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// CreateRecord creates a new analysis-artifact record
func (s *Service) CreateRecord(record *AnalysisArtifactRecord) error {
	// Validate required fields
	if record.RequestID == 0 {
		return fmt.Errorf("request_id is required")
	}
	if record.ArtifactID == 0 {
		return fmt.Errorf("artifact_id is required")
	}
	if record.Quantity <= 0 {
		return fmt.Errorf("quantity must be greater than 0")
	}

	return s.repo.CreateRecord(record)
}

// GetRecordByCompositeKey retrieves a record by its composite key
func (s *Service) GetRecordByCompositeKey(requestID, artifactID uint) (*AnalysisArtifactRecord, error) {
	return s.repo.GetRecordByCompositeKey(requestID, artifactID)
}

// GetRecordsByRequestID retrieves all records for a specific request
func (s *Service) GetRecordsByRequestID(requestID uint) ([]AnalysisArtifactRecord, error) {
	return s.repo.GetRecordsByRequestID(requestID)
}

// GetRecordsByArtifactID retrieves all records for a specific artifact
func (s *Service) GetRecordsByArtifactID(artifactID uint) ([]AnalysisArtifactRecord, error) {
	return s.repo.GetRecordsByArtifactID(artifactID)
}

// UpdateRecord updates specific fields of a record
func (s *Service) UpdateRecord(requestID, artifactID uint, updates map[string]interface{}) error {
	// Validate that the record exists
	_, err := s.repo.GetRecordByCompositeKey(requestID, artifactID)
	if err != nil {
		return fmt.Errorf("record not found: %w", err)
	}

	// Validate quantity if being updated
	if quantity, ok := updates["quantity"]; ok {
		if q, ok := quantity.(int); ok && q <= 0 {
			return fmt.Errorf("quantity must be greater than 0")
		}
	}

	return s.repo.UpdateRecord(requestID, artifactID, updates)
}

// DeleteRecord deletes a record from the request
func (s *Service) DeleteRecord(requestID, artifactID uint) error {
	// Validate that the record exists
	_, err := s.repo.GetRecordByCompositeKey(requestID, artifactID)
	if err != nil {
		return fmt.Errorf("record not found: %w", err)
	}

	return s.repo.DeleteRecord(requestID, artifactID)
}

// DeleteRecordsByRequestID deletes all records for a specific request
func (s *Service) DeleteRecordsByRequestID(requestID uint) error {
	return s.repo.DeleteRecordsByRequestID(requestID)
}