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

func (s *Service) CreateRecord(record *AnalysisArtifactRecord) error {
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

func (s *Service) GetRecordByCompositeKey(requestID, artifactID uint) (*AnalysisArtifactRecord, error) {
	return s.repo.GetRecordByCompositeKey(requestID, artifactID)
}

func (s *Service) GetRecordsByRequestID(requestID uint) ([]AnalysisArtifactRecord, error) {
	return s.repo.GetRecordsByRequestID(requestID)
}

func (s *Service) GetRecordsByArtifactID(artifactID uint) ([]AnalysisArtifactRecord, error) {
	return s.repo.GetRecordsByArtifactID(artifactID)
}

func (s *Service) UpdateRecord(requestID, artifactID uint, updates map[string]interface{}) error {
	_, err := s.repo.GetRecordByCompositeKey(requestID, artifactID)
	if err != nil {
		return fmt.Errorf("record not found: %w", err)
	}

	if quantity, ok := updates["quantity"]; ok {
		if q, ok := quantity.(int); ok && q <= 0 {
			return fmt.Errorf("quantity must be greater than 0")
		}
	}

	return s.repo.UpdateRecord(requestID, artifactID, updates)
}

func (s *Service) DeleteRecord(requestID, artifactID uint) error {
	_, err := s.repo.GetRecordByCompositeKey(requestID, artifactID)
	if err != nil {
		return fmt.Errorf("record not found: %w", err)
	}

	return s.repo.DeleteRecord(requestID, artifactID)
}

func (s *Service) DeleteRecordsByRequestID(requestID uint) error {
	return s.repo.DeleteRecordsByRequestID(requestID)
}