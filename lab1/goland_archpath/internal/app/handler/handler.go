package handler

import (
	"archpath/internal/app/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const activeCartID = "abc"

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) getCartStatus() (cartID string, analisedCount int) {
	cartID = activeCartID
	cartData, err := h.Repository.GetExcavationCart(cartID)

	if err != nil {
		logrus.Warnf("Could not fetch cart data %s for count: %v", cartID, err)
		return cartID, 0
	}

	if count, ok := cartData["TotalEntryCount"].(int); ok {
		analisedCount = count
	}
	return cartID, analisedCount
}

func (h *Handler) GetArtifactTypes(ctx *gin.Context) {
	searchQuery := ctx.Query("query")
	var commodities []repository.Artifact
	var err error

	if searchQuery == "" {
		commodities, err = h.Repository.GetCommodities()
	} else {
		commodities, err = h.Repository.GetCommoditiesByName(searchQuery)
	}

	if err != nil {
		logrus.Error("Error fetching artifact list:", err)
		commodities = []repository.Artifact{}
	}

	cartID, analisedCount := h.getCartStatus()

	ctx.HTML(http.StatusOK, "mainPage.html", gin.H{
		"commodities":   commodities,
		"query":         searchQuery,
		"analisedCount": analisedCount,
		"cartID":        cartID,
	})
}

func (h *Handler) GetArtifactTypeDetails(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error("Invalid artifact ID:", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	artifact, err := h.Repository.GetArtifact(id)
	if err != nil {
		logrus.Error("Error fetching artifact details:", err)
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	cartID, analisedCount := h.getCartStatus()

	ctx.HTML(http.StatusOK, "detailsPage.html", gin.H{
		"artifact":      artifact,
		"analisedCount": analisedCount,
		"cartID":        cartID,
	})
}

func (h *Handler) GetExcavationCart(ctx *gin.Context) {
	cartID := ctx.Param("id")

	cartData, err := h.Repository.GetExcavationCart(cartID)
	if err != nil {
		logrus.Error("Error fetching cart data:", err)
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.HTML(http.StatusOK, "cartPage.html", gin.H{
		"cart": cartData,
	})
}
