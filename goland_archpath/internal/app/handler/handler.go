package handler

import (
	"archpath/internal/app/repository"
	"archpath/internal/app/service"
)

const activeUserID = 1

// Handler structure uses the updated AnalysisService.
type Handler struct {
	Repository      *repository.Repository
	AnalysisService *service.AnalysisService
}

// NewHandler constructor uses the updated AnalysisService.
func NewHandler(r *repository.Repository, s *service.AnalysisService) *Handler {
	return &Handler{
		Repository:      r,
		AnalysisService: s,
	}
}

// getAnalysisStatus is a helper method used across handlers.
func (h *Handler) getAnalysisStatus() (analysisID uint, entryCount int) {
	return h.AnalysisService.GetAnalysisStatus(activeUserID)
}