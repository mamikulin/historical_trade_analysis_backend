package handler

import (
	"archpath/internal/app/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// GetTradeAnalysis retrieves and displays a trade analysis (by ID or draft for current user).
func (h *Handler) GetTradeAnalysis(ctx *gin.Context) {
	var analysis models.TradeAnalysis
	var err error

	analysisIDStr := ctx.Param("id")
	if analysisIDStr != "" {
		analysisID, err := strconv.ParseUint(analysisIDStr, 10, 32)
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		analysis, err = h.Repository.GetAnalysisByID(uint(analysisID))
		if err == gorm.ErrRecordNotFound {
			logrus.Warnf("Analysis with ID %d not found.", analysisID)
			ctx.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}
	} else {
		analysis, err = h.Repository.GetAnalysisByUser(activeUserID, "draft")
		if err == gorm.ErrRecordNotFound {
			logrus.Warnf("Draft analysis for user %d not found.", activeUserID)
			ctx.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}
	}

	if analysis.Status == "deleted" {
		ctx.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	if err != nil {
		logrus.Error("Error fetching analysis data:", err)
		ctx.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	totalItemCount := 0
	for _, entry := range analysis.Entries {
		totalItemCount += entry.Quantity
	}

	ctx.HTML(http.StatusOK, "cartPage.html", gin.H{
		"analysis":       analysis,
		"TotalItemCount": totalItemCount,
	})
}

// AddArtifactToAnalysis adds an artifact to the user's draft analysis.
func (h *Handler) AddArtifactToAnalysis(ctx *gin.Context) {
	artifactIDStr := ctx.PostForm("artifact_id")
	quantityStr := ctx.PostForm("quantity")

	artifactID, err := strconv.ParseUint(artifactIDStr, 10, 32)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	quantity, err := strconv.Atoi(quantityStr)
	if err != nil || quantity <= 0 {
		quantity = 1
	}

	analysis, err := h.AnalysisService.AddArtifactToAnalysis(activeUserID, uint(artifactID), quantity)
	if err != nil {
		logrus.Errorf("Error adding artifact via service: %v", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Redirect(http.StatusFound, "/analysis/"+strconv.FormatUint(uint64(analysis.ID), 10))
}

// DeleteTradeAnalysis soft deletes an analysis record by setting its status to "deleted".
func (h *Handler) DeleteTradeAnalysis(ctx *gin.Context) {
	analysisIDStr := ctx.PostForm("analysis_id")
	analysisID, err := strconv.ParseUint(analysisIDStr, 10, 32)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = h.Repository.UpdateAnalysisStatusSQL(uint(analysisID), "deleted")
	if err != nil {
		logrus.Errorf("Error soft deleting analysis %d: %v", analysisID, err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Redirect(http.StatusFound, "/")
}

// RemoveArtifactFromAnalysis removes a specific artifact record from an analysis record.
func (h *Handler) RemoveArtifactFromAnalysis(ctx *gin.Context) {
	analysisIDStr := ctx.PostForm("analysis_id")
	artifactIDStr := ctx.PostForm("artifact_id")

	analysisID, err := strconv.ParseUint(analysisIDStr, 10, 32)
	if err != nil {
		logrus.Errorf("Invalid analysis_id: %v", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	artifactID, err := strconv.ParseUint(artifactIDStr, 10, 32)
	if err != nil {
		logrus.Errorf("Invalid artifact_id: %v", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = h.Repository.RemoveArtifactFromAnalysis(uint(analysisID), uint(artifactID))
	if err != nil {
		logrus.Errorf("Error deleting analysis record A:%d I:%d: %v", analysisID, artifactID, err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Redirect(http.StatusFound, "/analysis/"+analysisIDStr)
}

// UpdateTradeAnalysis handles updating details (like the SiteName) of an analysis record.
func (h *Handler) UpdateTradeAnalysis(ctx *gin.Context) {
	analysisIDStr := ctx.Param("id")
	analysisID, err := strconv.ParseUint(analysisIDStr, 10, 32)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	siteName := ctx.PostForm("site_name")

	if siteName == "" {
		logrus.Error("Site name cannot be empty during analysis update.")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = h.Repository.UpdateAnalysisDetails(uint(analysisID), siteName)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logrus.Warnf("Analysis with ID %d not found for update.", analysisID)
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}
		logrus.Errorf("Error updating analysis details for ID %d: %v", analysisID, err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Redirect(http.StatusFound, "/analysis/"+analysisIDStr)
}

// UpdateArtifactQuantityInAnalysis handles changing the quantity of a specific artifact record.
func (h *Handler) UpdateArtifactQuantityInAnalysis(ctx *gin.Context) {
	analysisIDStr := ctx.PostForm("analysis_id")
	artifactIDStr := ctx.PostForm("artifact_id")
	quantityStr := ctx.PostForm("quantity")

	analysisID, err := strconv.ParseUint(analysisIDStr, 10, 32)
	if err != nil {
		logrus.Errorf("Invalid analysis_id: %v", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	artifactID, err := strconv.ParseUint(artifactIDStr, 10, 32)
	if err != nil {
		logrus.Errorf("Invalid artifact_id: %v", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	quantity, err := strconv.Atoi(quantityStr)
	if err != nil || quantity <= 0 {
		logrus.Errorf("Invalid or non-positive quantity: %v", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = h.Repository.UpdateArtifactQuantityInAnalysis(uint(analysisID), uint(artifactID), quantity)
	if err != nil {
		logrus.Errorf("Error updating quantity A:%d I:%d Q:%d: %v", analysisID, artifactID, quantity, err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Redirect(http.StatusFound, "/analysis/"+analysisIDStr)
}