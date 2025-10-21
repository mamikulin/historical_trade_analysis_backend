package service

import (
	"archpath/internal/app/models"
	"archpath/internal/app/repository"
	"fmt"

	"gorm.io/gorm"
)

type AnalysisService struct {
	repo *repository.Repository
}

func NewAnalysisService(r *repository.Repository) *AnalysisService {
	return &AnalysisService{
		repo: r,
	}
}

// GetOrCreateDraftAnalysis checks for a draft analysis for the user, or creates one if not found.
func (s *AnalysisService) GetOrCreateDraftAnalysis(userID uint) (models.TradeAnalysis, error) {
    // Note: requires import "fmt" for Sprintf

    const draftStatus = "draft"
    analysis, err := s.repo.GetAnalysisByUser(userID, draftStatus)
    
    if err != nil && err == gorm.ErrRecordNotFound {
        // 1. Create new analysis with a temporary SiteName
        newAnalysis := models.TradeAnalysis{
            Status:    draftStatus,
            CreatorID: userID,
            SiteName:  "Черновик", // Временное имя
        }
        if err := s.repo.CreateAnalysis(&newAnalysis); err != nil {
            return models.TradeAnalysis{}, fmt.Errorf("failed to create analysis: %w", err)
        }
        
        // 2. Update SiteName using the newly generated ID
        newAnalysis.SiteName = fmt.Sprintf("Памятник №%d", newAnalysis.ID)
        
        // 3. Save the updated SiteName back to the repository
        // Предполагается наличие метода UpdateAnalysis в репозитории.
        if err := s.repo.UpdateAnalysis(&newAnalysis); err != nil {
            return models.TradeAnalysis{}, fmt.Errorf("failed to update analysis SiteName: %w", err)
        }
        
        // 4. Return the newly created analysis (need to fetch it to populate Entries)
        return s.repo.GetAnalysisByID(newAnalysis.ID)
    } else if err != nil {
        return models.TradeAnalysis{}, fmt.Errorf("analysis search error: %w", err)
    }
    
    return analysis, nil
}

// AddArtifactToAnalysis handles finding/creating the analysis and then adding/updating the entry.
func (s *AnalysisService) AddArtifactToAnalysis(userID uint, artifactID uint, quantity int) (models.TradeAnalysis, error) {
	analysis, err := s.GetOrCreateDraftAnalysis(userID)
	if err != nil {
		return models.TradeAnalysis{}, err
	}

	record, err := s.repo.GetAnalysisRecord(analysis.ID, artifactID)
	
	if err != nil && err == gorm.ErrRecordNotFound {
		// Create new record
		record = models.AnalysisArtifactRecord{
			RequestID:  analysis.ID,
			ArtifactID: artifactID,
			Quantity:   quantity,
		}
		if err := s.repo.CreateAnalysisRecord(&record); err != nil {
			return models.TradeAnalysis{}, fmt.Errorf("failed to add artifact to analysis: %w", err)
		}
	} else if err != nil {
		return models.TradeAnalysis{}, fmt.Errorf("record search error: %w", err)
	} else {
		// Update existing record - increment quantity
		record.Quantity += quantity
		if err := s.repo.UpdateAnalysisRecord(&record); err != nil {
			return models.TradeAnalysis{}, fmt.Errorf("failed to update artifact in analysis: %w", err)
		}
	}

	// Return the updated analysis with preloaded entries
	return s.repo.GetAnalysisByID(analysis.ID)
}

// GetAnalysisStatus is a simple helper function for the handler to check analysis state.
func (s *AnalysisService) GetAnalysisStatus(userID uint) (analysisID uint, entryCount int) {
	analysis, err := s.repo.GetAnalysisByUser(userID, "draft")
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			// In a real app, logrus.Errorf("Error searching for draft analysis: %v", err)
		}
		return 0, 0
	}

	count := 0
	for _, entry := range analysis.Entries {
		count += entry.Quantity
	}

	return analysis.ID, count
}