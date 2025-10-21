package handler

import (
	"archpath/internal/app/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GetArtifactTypes handles the listing of artifacts, including search filtering.
func (h *Handler) GetArtifactTypes(ctx *gin.Context) {
	searchQuery := ctx.Query("query")
	artifacts, err := h.Repository.GetArtifacts(searchQuery)

	if err != nil {
		logrus.Error("Error fetching artifact list:", err)
		artifacts = []models.Artifact{}
	}

	analysisID, entryCount := h.getAnalysisStatus()

	ctx.HTML(http.StatusOK, "mainPage.html", gin.H{
		"commodities":   artifacts,
		"query":         searchQuery,
		"analisedCount": entryCount,
		"analysisID":    analysisID,
	})
}

// GetArtifactTypeDetails handles fetching the details for a single artifact.
func (h *Handler) GetArtifactTypeDetails(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	artifact, err := h.Repository.GetArtifactByID(uint(id))
	if err != nil {
		logrus.Error("Error fetching artifact details:", err)
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	analysisID, entryCount := h.getAnalysisStatus()

	ctx.HTML(http.StatusOK, "detailsPage.html", gin.H{
		"artifact":      artifact,
		"analisedCount": entryCount,
		"analysisID":    analysisID,
	})
}