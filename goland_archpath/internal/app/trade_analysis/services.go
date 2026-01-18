package trade_analysis

import (
	"archpath/internal/app/analysis_artifact_record"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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

func (s *Service) GetAllRequests(status string, startDate, endDate *time.Time, creatorID *uint) ([]map[string]interface{}, error) {
	requests, err := s.repo.GetAllRequests(status, startDate, endDate, creatorID)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(requests))
	for i, req := range requests {
		// Подсчитываем количество записей с calculated_value
		calculatedCount, err := s.repo.CountCompletedEntries(req.ID)
		if err != nil {
			calculatedCount = 0
		}

		result[i] = map[string]interface{}{
			"id":                       req.ID,
			"status":                   req.Status,
			"creator_id":               req.CreatorID,
			"site_name":                req.SiteName,
			"formation_date":           req.FormationDate,
			"completion_date":          req.CompletionDate,
			"moderator_id":             req.ModeratorID,
			"calculated_entries_count": calculatedCount,
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

	// Подсчитываем количество записей с calculated_value
	calculatedCount, err := s.repo.CountCompletedEntries(id)
	if err != nil {
		calculatedCount = 0
	}

	return map[string]interface{}{
		"id":                       request.ID,
		"status":                   request.Status,
		"creator_id":               request.CreatorID,
		"site_name":                request.SiteName,
		"formation_date":           request.FormationDate,
		"completion_date":          request.CompletionDate,
		"moderator_id":             request.ModeratorID,
		"entries":                  entries,
		"calculated_entries_count": calculatedCount,
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
		log.Printf("Заявка ID: %d одобрена модератором ID: %d", id, moderatorID)
		entries, err := s.repo.GetEntriesWithArtifacts(id)
		if err != nil {
			return fmt.Errorf("failed to retrieve entries: %w", err)
		}

		// Запускаем async-расчет для каждой записи (НЕ ждем результата)
		go s.triggerAsyncCalculations(id, entries)
	} else {
		log.Printf("Заявка ID: %d отклонена модератором ID: %d", id, moderatorID)
	}

	return s.repo.UpdateRequest(id, updates)
}



// triggerAsyncCalculations вызывает Django-сервис для расчета всех м-м записей
// Все вычисления процентов выполняются в Django
func (s *Service) triggerAsyncCalculations(requestID uint, entries []analysis_artifact_record.AnalysisArtifactRecord) {
	asyncServiceURL := "http://localhost:8001/api/calculate" // URL Django-сервиса
	
	log.Printf("Начат расчет для заявки ID: %d (%d артефактов)", requestID, len(entries))
	
	// Собираем данные о всех записях для отправки в Django
	type EntryData struct {
		ArtifactID       uint   `json:"artifact_id"`
		ProductionCenter string `json:"production_center"`
		Quantity         int    `json:"quantity"`
	}
	
	var entryDataList []EntryData
	for _, entry := range entries {
		entryDataList = append(entryDataList, EntryData{
			ArtifactID:       entry.ArtifactID,
			ProductionCenter: entry.Artifact.ProductionCenter,
			Quantity:         entry.Quantity,
		})
	}
	
	// Формируем шаблон callback URL (Django заменит {request_id} и {artifact_id})
	callbackURLTemplate := "http://localhost:8000/api/trade-analysis/{request_id}/entries/{artifact_id}/result"
	
	payload := map[string]interface{}{
		"request_id":   requestID,
		"entries":      entryDataList,
		"callback_url": callbackURLTemplate,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal payload for async calculation: %v", err)
		return
	}

	log.Printf("Отправка данных в Django для обработки...")
	
	// HTTP POST к Django (не ждем ответа)
	go func(data []byte) {
		resp, err := http.Post(asyncServiceURL, "application/json", bytes.NewBuffer(data))
		if err != nil {
			log.Printf("Failed to call async service: %v", err)
			return
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			log.Printf("Async service returned non-OK status: %d, body: %s", resp.StatusCode, string(bodyBytes))
		}
	}(jsonData)
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