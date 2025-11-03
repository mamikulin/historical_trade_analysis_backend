package trade_analysis

import (
	"fmt"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetDraftCart(creatorID uint) (map[string]interface{}, error) {
	draft, err := s.repo.GetDraftByCreatorID(creatorID)
	if err != nil {
		draft = &TradeAnalysis{
			Status:    "draft",
			CreatorID: creatorID,
			SiteName:  "",
		}
		if err := s.repo.CreateRequest(draft); err != nil {
			return nil, fmt.Errorf("failed to create draft: %w", err)
		}
	}
	
	count, err := s.repo.CountEntriesByRequestID(draft.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to count entries: %w", err)
	}
	
	return map[string]interface{}{
		"request_id":    draft.ID,
		"entries_count": count,
	}, nil
}

func (s *Service) GetAllRequests(status string, startDate, endDate *time.Time) ([]map[string]interface{}, error) {
	requests, err := s.repo.GetAllRequests(status, startDate, endDate)
	if err != nil {
		return nil, err
	}
	
	result := make([]map[string]interface{}, len(requests))
	for i, req := range requests {
		result[i] = map[string]interface{}{
			"id":                   req.ID,
			"status":               req.Status,
			"creator_id":           req.CreatorID,
			"site_name":            req.SiteName,
			"formation_date":       req.FormationDate,
			"completion_date":      req.CompletionDate,
			"moderator_id":         req.ModeratorID,
			"total_finds_quantity": req.TotalFindsQuantity,
			"analysis_result":      req.AnalysisResult,
		}
	}
	
	return result, nil
}

func (s *Service) GetRequestByID(id uint) (map[string]interface{}, error) {
	request, err := s.repo.GetRequestByID(id)
	if err != nil {
		return nil, fmt.Errorf("request not found: %w", err)
	}
	
	entries, err := s.repo.GetEntriesWithArtifacts(id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve entries: %w", err)
	}
	
	return map[string]interface{}{
		"id":                   request.ID,
		"status":               request.Status,
		"creator_id":           request.CreatorID,
		"site_name":            request.SiteName,
		"formation_date":       request.FormationDate,
		"completion_date":      request.CompletionDate,
		"moderator_id":         request.ModeratorID,
		"total_finds_quantity": request.TotalFindsQuantity,
		"analysis_result":      request.AnalysisResult,
		"entries":              entries,
	}, nil
}

func (s *Service) UpdateRequest(id uint, updates map[string]interface{}) error {
	_, err := s.repo.GetRequestByID(id)
	if err != nil {
		return fmt.Errorf("request not found: %w", err)
	}
	
	return s.repo.UpdateRequest(id, updates)
}

func (s *Service) FormRequest(id uint, creatorID uint) error {
	request, err := s.repo.GetRequestByID(id)
	if err != nil {
		return fmt.Errorf("request not found: %w", err)
	}
	
	if request.CreatorID != creatorID {
		return fmt.Errorf("unauthorized: only creator can form the request")
	}
	
	if request.SiteName == "" {
		return fmt.Errorf("site_name is required")
	}
	
	count, err := s.repo.CountEntriesByRequestID(id)
	if err != nil {
		return fmt.Errorf("failed to count entries: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("cannot form request without entries")
	}
	
	now := time.Now()
	updates := map[string]interface{}{
		"status":         "formed",
		"formation_date": &now,
	}
	
	return s.repo.UpdateRequest(id, updates)
}

func (s *Service) CompleteOrRejectRequest(id uint, moderatorID uint, action string) error {
	request, err := s.repo.GetRequestByID(id)
	if err != nil {
		return fmt.Errorf("request not found: %w", err)
	}
	
	if action != "completed" && action != "rejected" {
		return fmt.Errorf("invalid action: must be 'completed' or 'rejected'")
	}
	
	if request.Status != "formed" {
		return fmt.Errorf("can only complete/reject formed requests")
	}
	
	now := time.Now()
	updates := map[string]interface{}{
		"status":          action,
		"moderator_id":    moderatorID,
		"completion_date": &now,
	}
	
	if action == "completed" {
		entries, err := s.repo.GetEntriesWithArtifacts(id)
		if err != nil {
			return fmt.Errorf("failed to retrieve entries: %w", err)
		}
		
		analysisResult := request.GetPercentageByRegion(entries)
		updates["analysis_result"] = analysisResult
	}
	
	return s.repo.UpdateRequest(id, updates)
}

func (s *Service) DeleteRequest(id uint) error {
	request, err := s.repo.GetRequestByID(id)
	if err != nil {
		return fmt.Errorf("request not found: %w", err)
	}
	
	if request.FormationDate == nil {
		return fmt.Errorf("can only delete formed requests")
	}
	
	return s.repo.DeleteRequest(id)
}