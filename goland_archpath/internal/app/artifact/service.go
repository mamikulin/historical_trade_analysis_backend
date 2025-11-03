package artifact

import (
	"fmt"
	"mime/multipart"
)

type Service interface {
	GetAll(filters map[string]interface{}) ([]Artifact, error)
	GetByID(id uint) (*Artifact, error)
	Create(a *Artifact) error
	Update(id uint, data Artifact) error
	Delete(id uint) error
	UploadImage(id uint, file multipart.File, header *multipart.FileHeader) (string, error)
	AddToDraft(artifactID, creatorID uint, quantity int, comment string) (map[string]interface{}, error)
}

type MinioClient interface {
	UploadFile(objectName string, file multipart.File, header *multipart.FileHeader) (string, error)
}

type TradeAnalysisService interface {
	GetDraftCart(creatorID uint) (map[string]interface{}, error)
}

type AnalysisArtifactRecordRepository interface {
	CreateRecord(record interface{}) error
	GetRecordByCompositeKey(requestID, artifactID uint) (interface{}, error)
	UpdateRecord(requestID, artifactID uint, updates map[string]interface{}) error
}

type service struct {
	repo        *Repository
	minioClient MinioClient
	taService   TradeAnalysisService
	aarRepo     AnalysisArtifactRecordRepository
}

func NewService(repo *Repository, minioClient MinioClient) Service {
	return &service{repo: repo, minioClient: minioClient}
}

// SetTradeAnalysisService sets the trade analysis service (for avoiding circular dependency)
func (s *service) SetTradeAnalysisService(taService TradeAnalysisService) {
	s.taService = taService
}

// SetAnalysisArtifactRecordRepository sets the AAR repository (for avoiding circular dependency)
func (s *service) SetAnalysisArtifactRecordRepository(aarRepo AnalysisArtifactRecordRepository) {
	s.aarRepo = aarRepo
}

func (s *service) GetAll(filters map[string]interface{}) ([]Artifact, error) {
	return s.repo.GetAll()
}

func (s *service) GetByID(id uint) (*Artifact, error) {
	return s.repo.GetByID(id)
}

func (s *service) Create(a *Artifact) error {
	return s.repo.Create(a)
}

func (s *service) Update(id uint, data Artifact) error {
	return s.repo.Update(id, data)
}

func (s *service) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *service) UploadImage(id uint, file multipart.File, header *multipart.FileHeader) (string, error) {
	url, err := s.minioClient.UploadFile(fmt.Sprintf("artifact_%d", id), file, header)
	if err != nil {
		return "", err
	}
	artifact, _ := s.repo.GetByID(id)
	artifact.ImageURL = &url
	_ = s.repo.Update(id, *artifact)
	return url, nil
}

// AddToDraft adds an artifact to the user's draft request
func (s *service) AddToDraft(artifactID, creatorID uint, quantity int, comment string) (map[string]interface{}, error) {
	if s.taService == nil || s.aarRepo == nil {
		return nil, fmt.Errorf("dependencies not set")
	}
	
	// Verify artifact exists
	artifact, err := s.repo.GetByID(artifactID)
	if err != nil {
		return nil, fmt.Errorf("artifact not found: %w", err)
	}
	
	// Get or create draft request
	cart, err := s.taService.GetDraftCart(creatorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get draft: %w", err)
	}
	
	requestID := cart["request_id"].(uint)
	
	// Check if artifact already in draft
	existing, err := s.aarRepo.GetRecordByCompositeKey(requestID, artifactID)
	if err == nil && existing != nil {
		// Update existing record
		updates := map[string]interface{}{
			"quantity": quantity,
			"comment":  comment,
		}
		if err := s.aarRepo.UpdateRecord(requestID, artifactID, updates); err != nil {
			return nil, fmt.Errorf("failed to update existing record: %w", err)
		}
		
		return map[string]interface{}{
			"message":     "Artifact quantity updated in draft",
			"request_id":  requestID,
			"artifact_id": artifactID,
			"quantity":    quantity,
		}, nil
	}
	
	// Create new record
	record := map[string]interface{}{
		"request_id":  requestID,
		"artifact_id": artifactID,
		"quantity":    quantity,
		"comment":     comment,
		"order":       0,
	}
	
	if err := s.aarRepo.CreateRecord(record); err != nil {
		return nil, fmt.Errorf("failed to add to draft: %w", err)
	}
	
	return map[string]interface{}{
		"message":      "Artifact added to draft",
		"request_id":   requestID,
		"artifact_id":  artifactID,
		"artifact_name": artifact.Name,
		"quantity":     quantity,
	}, nil
}