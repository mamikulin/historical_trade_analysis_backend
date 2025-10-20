package service

import (
	"archpath/internal/app/models"
	"archpath/internal/app/repository"
	"fmt"

	"gorm.io/gorm"
)

type CartService struct {
	repo *repository.Repository
}

func NewCartService(r *repository.Repository) *CartService {
	return &CartService{
		repo: r,
	}
}

// GetOrCreateDraftCart checks for a draft cart for the user, or creates one if not found.
func (s *CartService) GetOrCreateDraftCart(userID uint) (models.SiteCart, error) {
	const draftStatus = "draft"
	cart, err := s.repo.GetSiteCartByUser(userID, draftStatus)

	if err != nil && err == gorm.ErrRecordNotFound {
		// Create new cart
		newCart := models.SiteCart{
			Status:    draftStatus,
			CreatorID: userID,
			SiteName:  "New site", // Default site name
		}
		if err := s.repo.CreateSiteCart(&newCart); err != nil {
			return models.SiteCart{}, fmt.Errorf("failed to create cart: %w", err)
		}
		// Return the newly created cart (need to fetch it to populate Entries)
		return s.repo.GetSiteCartByID(newCart.ID)

	} else if err != nil {
		return models.SiteCart{}, fmt.Errorf("cart search error: %w", err)
	}

	return cart, nil
}

// AddArtifactToCart handles finding/creating the cart and then adding/updating the entry.
func (s *CartService) AddArtifactToCart(userID uint, artifactID uint, quantity int) (models.SiteCart, error) {
	cart, err := s.GetOrCreateDraftCart(userID)
	if err != nil {
		return models.SiteCart{}, err
	}

	entry, err := s.repo.GetSiteEntry(cart.ID, artifactID)

	if err != nil && err == gorm.ErrRecordNotFound {
		// Create new entry
		entry = models.SiteEntry{
			CartID:           cart.ID,
			ArtifactID:       artifactID,
			ArtifactQuantity: quantity,
		}
		if err := s.repo.CreateSiteEntry(&entry); err != nil {
			return models.SiteCart{}, fmt.Errorf("failed to add artifact to cart: %w", err)
		}
	} else if err != nil {
		return models.SiteCart{}, fmt.Errorf("entry search error: %w", err)
	} else {
		// Update existing entry
		entry.ArtifactQuantity += quantity
		if err := s.repo.UpdateSiteEntry(&entry); err != nil {
			return models.SiteCart{}, fmt.Errorf("failed to update artifact in cart: %w", err)
		}
	}

	// Return the updated cart with preloaded entries
	return s.repo.GetSiteCartByID(cart.ID)
}

// GetCartStatus is a simple helper function for the handler to check cart state.
func (s *CartService) GetCartStatus(userID uint) (cartID uint, entryCount int) {
	cart, err := s.repo.GetSiteCartByUser(userID, "draft")
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			// In a real app, logrus.Errorf("Error searching for draft cart: %v", err)
		}
		return 0, 0
	}

	count := 0
	for _, entry := range cart.Entries {
		count += entry.ArtifactQuantity
	}

	return cart.ID, count
}