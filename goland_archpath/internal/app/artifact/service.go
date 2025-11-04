package artifact

import (
	"fmt"
	"mime/multipart"
	"time"
	
	"gorm.io/gorm"
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

type service struct {
	repo        *Repository
	minioClient MinioClient
}

func NewService(repo *Repository, minioClient MinioClient) Service {
	return &service{
		repo:        repo,
		minioClient: minioClient,
	}
}

func (s *service) GetAll(filters map[string]interface{}) ([]Artifact, error) {
	return s.repo.GetAll(filters)
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

// TradeAnalysisDraft represents a minimal draft structure
type TradeAnalysisDraft struct {
	ID        uint
	Status    string
	CreatorID uint
	SiteName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// AnalysisArtifactRecord represents the many-to-many record
type AnalysisArtifactRecord struct {
	RequestID   uint
	ArtifactID  uint
	Quantity    int
	Comment     string
	Order       int
	IsMainEntry bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// AddToDraft adds an artifact to the user's draft request
// Заявка создается автоматически с указанием создателя, даты создания и статуса
func (s *service) AddToDraft(artifactID, creatorID uint, quantity int, comment string) (map[string]interface{}, error) {
	// Verify artifact exists
	artifact, err := s.repo.GetByID(artifactID)
	if err != nil {
		return nil, fmt.Errorf("artifact not found: %w", err)
	}
	
	// Get or create draft request
	var draft TradeAnalysisDraft
	err = s.repo.DB.Table("trade_analyses").
		Where("creator_id = ? AND status = ?", creatorID, "draft").
		First(&draft).Error
	
	if err == gorm.ErrRecordNotFound {
		// Create new draft - заявка создается пустой
		draft = TradeAnalysisDraft{
			Status:    "draft",
			CreatorID: creatorID,
			SiteName:  "",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		
		if err := s.repo.DB.Table("trade_analyses").Create(&draft).Error; err != nil {
			return nil, fmt.Errorf("failed to create draft: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	
	requestID := draft.ID
	
	// Check if artifact already in draft
	var existing AnalysisArtifactRecord
	err = s.repo.DB.Table("analysis_artifact_records").
		Where("request_id = ? AND artifact_id = ?", requestID, artifactID).
		First(&existing).Error
	
	if err == nil {
		// Update existing record
		if err := s.repo.DB.Table("analysis_artifact_records").
			Where("request_id = ? AND artifact_id = ?", requestID, artifactID).
			Updates(map[string]interface{}{
				"quantity":   quantity,
				"comment":    comment,
				"updated_at": time.Now(),
			}).Error; err != nil {
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
	record := AnalysisArtifactRecord{
		RequestID:   requestID,
		ArtifactID:  artifactID,
		Quantity:    quantity,
		Comment:     comment,
		Order:       0,
		IsMainEntry: false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	if err := s.repo.DB.Table("analysis_artifact_records").Create(&record).Error; err != nil {
		return nil, fmt.Errorf("failed to add to draft: %w", err)
	}
	
	return map[string]interface{}{
		"message":       "Artifact added to draft",
		"request_id":    requestID,
		"artifact_id":   artifactID,
		"artifact_name": artifact.Name,
		"quantity":      quantity,
	}, nil
}