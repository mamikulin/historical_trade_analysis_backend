package handler

import (
	"archpath/internal/app/models"
	"archpath/internal/app/repository"
	"archpath/internal/app/service" // New dependency
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const activeUserID = 1

type Handler struct {
	Repository  *repository.Repository // For artifact-related reads
	CartService *service.CartService   // For cart-related logic
}

func NewHandler(r *repository.Repository, s *service.CartService) *Handler {
	return &Handler{
		Repository:  r,
		CartService: s,
	}
}

// getCartStatus delegates to the service
func (h *Handler) getCartStatus() (cartID uint, entryCount int) {
	return h.CartService.GetCartStatus(activeUserID)
}

func (h *Handler) GetArtifactTypes(ctx *gin.Context) {
	searchQuery := ctx.Query("query")
	artifacts, err := h.Repository.GetArtifacts(searchQuery)

	if err != nil {
		logrus.Error("Error fetching artifact list:", err)
		artifacts = []models.Artifact{}
	}

	cartID, entryCount := h.getCartStatus() // Uses service layer

	ctx.HTML(http.StatusOK, "mainPage.html", gin.H{
		"commodities":   artifacts,
		"query":         searchQuery,
		"analisedCount": entryCount,
		"cartID":        cartID,
	})
}

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

	cartID, entryCount := h.getCartStatus() // Uses service layer

	ctx.HTML(http.StatusOK, "detailsPage.html", gin.H{
		"artifact":      artifact,
		"analisedCount": entryCount,
		"cartID":        cartID,
	})
}

func (h *Handler) GetSiteCart(ctx *gin.Context) {
	var cart models.SiteCart
	var err error

	cartIDStr := ctx.Param("id")
	if cartIDStr != "" {
		cartID, err := strconv.ParseUint(cartIDStr, 10, 32)
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		cart, err = h.Repository.GetSiteCartByID(uint(cartID))
		if err == gorm.ErrRecordNotFound {
			logrus.Warnf("Cart with ID %d not found.", cartID)
			ctx.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}
	} else {
		// Use repository directly for read operations
		cart, err = h.Repository.GetSiteCartByUser(activeUserID, "draft")
		if err == gorm.ErrRecordNotFound {
			logrus.Warnf("Draft cart for user %d not found.", activeUserID)
			ctx.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}
	}

	if cart.Status == "deleted" {
		ctx.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	if err != nil {
		logrus.Error("Error fetching cart data:", err)
		ctx.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	totalItemCount := 0
	for _, entry := range cart.Entries {
		totalItemCount += entry.ArtifactQuantity
	}

	ctx.HTML(http.StatusOK, "cartPage.html", gin.H{
		"cart":           cart,
		"TotalItemCount": totalItemCount,
	})
}

func (h *Handler) AddArtifactToCart(ctx *gin.Context) {
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

	// Logic is delegated entirely to the service layer
	cart, err := h.CartService.AddArtifactToCart(activeUserID, uint(artifactID), quantity)
	if err != nil {
		logrus.Errorf("Error adding artifact via service: %v", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Redirect(http.StatusFound, "/cart/"+strconv.FormatUint(uint64(cart.ID), 10))
}

func (h *Handler) DeleteSiteCart(ctx *gin.Context) {
	cartIDStr := ctx.PostForm("cart_id")
	cartID, err := strconv.ParseUint(cartIDStr, 10, 32)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = h.Repository.DeleteSiteCartSQL(uint(cartID), "deleted")
	if err != nil {
		logrus.Errorf("Error soft deleting cart %d: %v", cartID, err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Redirect(http.StatusFound, "/")
}

func (h *Handler) RemoveArtifactFromCart(ctx *gin.Context) {
	cartIDStr := ctx.PostForm("cart_id")
	artifactIDStr := ctx.PostForm("artifact_id")

	cartID, err := strconv.ParseUint(cartIDStr, 10, 32)
	if err != nil {
		logrus.Errorf("Invalid cart_id: %v", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	artifactID, err := strconv.ParseUint(artifactIDStr, 10, 32)
	if err != nil {
		logrus.Errorf("Invalid artifact_id: %v", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = h.Repository.RemoveArtifactFromCart(uint(cartID), uint(artifactID))
	if err != nil {
		logrus.Errorf("Error deleting cart entry C:%d A:%d: %v", cartID, artifactID, err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Redirect(http.StatusFound, "/cart/"+cartIDStr)
}

func (h *Handler) UpdateSiteCart(ctx *gin.Context) {
	cartIDStr := ctx.Param("id")
	cartID, err := strconv.ParseUint(cartIDStr, 10, 32)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	siteName := ctx.PostForm("site_name")
	comment := ctx.PostForm("comment")

	if siteName == "" {
		logrus.Error("Site name cannot be empty during cart update.")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = h.Repository.UpdateCartDetails(uint(cartID), siteName, comment)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logrus.Warnf("Cart with ID %d not found for update.", cartID)
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}
		logrus.Errorf("Error updating cart details for ID %d: %v", cartID, err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Redirect(http.StatusFound, "/cart/"+cartIDStr)
}

func (h *Handler) UpdateArtifactQuantityInCart(ctx *gin.Context) {
	cartIDStr := ctx.PostForm("cart_id")
	artifactIDStr := ctx.PostForm("artifact_id")
	quantityStr := ctx.PostForm("quantity")

	cartID, err := strconv.ParseUint(cartIDStr, 10, 32)
	if err != nil {
		logrus.Errorf("Invalid cart_id: %v", err)
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

	err = h.Repository.UpdateArtifactQuantityInCart(uint(cartID), uint(artifactID), quantity)
	if err != nil {
		logrus.Errorf("Error updating quantity C:%d A:%d Q:%d: %v", cartID, artifactID, quantity, err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Redirect(http.StatusFound, "/cart/"+cartIDStr)
}