package handler

import (
	"archpath/internal/app/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) GetCommodities(ctx *gin.Context) {
	var commodities []repository.Commodity
	var err error

	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		commodities, err = h.Repository.GetCommodities()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		commodities, err = h.Repository.GetCommoditiesByName(searchQuery)
		if err != nil {
			logrus.Error(err)
		}
	}

	analysedCommodities := h.getAnalisedCommodities()
	analisedCount := len(analysedCommodities)

	ctx.HTML(http.StatusOK, "listCommodities.html", gin.H{
		"time":          time.Now().Format("15:04:05"),
		"commodities":   commodities,
		"query":         searchQuery,
		"analisedCount": analisedCount,
	})
}

func (h *Handler) GetAnalysisPage(ctx *gin.Context) {
	analysedCommodities := h.getAnalisedCommodities()
	analisedCount := len(analysedCommodities)

	ctx.HTML(http.StatusOK, "analysis.html", gin.H{
		"time":             time.Now().Format("15:04:05"),
		"commoditiesCount": analisedCount,
		"commodities":      analysedCommodities,
	})
}

func (h *Handler) GetCommodity(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	commodity, err := h.Repository.GetCommodity(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "commodity.html", gin.H{
		"commodity": commodity,
	})
}

func (h *Handler) getAnalisedCommodities() []repository.Commodity {
	allCommodities, err := h.Repository.GetCommodities()
	if err != nil {
		return []repository.Commodity{}
	}
	return allCommodities[:4]
}

func (h *Handler) getAnalisedCount(ctx *gin.Context) {

	analisedCommodities := h.getAnalisedCommodities()
	analisedCount := len(analisedCommodities)

	ctx.JSON(http.StatusOK, gin.H{
		"count": analisedCount,
	})
}
