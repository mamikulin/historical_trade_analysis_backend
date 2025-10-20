package repository

import (
	"archpath/internal/app/models"
	"database/sql"
	"fmt"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	DB    *gorm.DB
	SQLDB *sql.DB
}

func NewRepository(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("database connection error: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("error getting *sql.DB: %w", err)
	}

	err = db.AutoMigrate(&models.User{}, &models.Artifact{}, &models.SiteCart{}, &models.SiteEntry{})
	if err != nil {
		return nil, fmt.Errorf("migration error: %w", err)
	}

	return &Repository{DB: db, SQLDB: sqlDB}, nil
}

func (r *Repository) GetArtifacts(searchQuery string) ([]models.Artifact, error) {
	var artifacts []models.Artifact
	query := r.DB.Where("is_active = ?", true)

	if searchQuery != "" {
		query = query.Where("LOWER(name) LIKE LOWER(?)", "%"+strings.ToLower(searchQuery)+"%")
	}

	if err := query.Find(&artifacts).Error; err != nil {
		return nil, err
	}
	return artifacts, nil
}

func (r *Repository) GetArtifactByID(id uint) (models.Artifact, error) {
	var artifact models.Artifact
	if err := r.DB.First(&artifact, id).Error; err != nil {
		return models.Artifact{}, err
	}
	return artifact, nil
}

func (r *Repository) GetSiteCartByUser(userID uint, status string) (models.SiteCart, error) {
	var cart models.SiteCart
	err := r.DB.
		Preload("Entries.Artifact").
		Preload("Creator").
		Where("creator_id = ? AND status = ?", userID, status).
		First(&cart).Error

	if err != nil {
		return models.SiteCart{}, err
	}
	return cart, nil
}

func (r *Repository) GetSiteCartByID(cartID uint) (models.SiteCart, error) {
	var cart models.SiteCart
	err := r.DB.
		Preload("Entries.Artifact").
		Preload("Creator").
		First(&cart, cartID).Error

	if err != nil {
		return models.SiteCart{}, err
	}
	return cart, nil
}

// --- NEW/SIMPLIFIED CART/ENTRY CRUD METHODS ---

func (r *Repository) CreateSiteCart(cart *models.SiteCart) error {
	return r.DB.Create(cart).Error
}

func (r *Repository) GetSiteEntry(cartID uint, artifactID uint) (models.SiteEntry, error) {
	var entry models.SiteEntry
	err := r.DB.Where("cart_id = ? AND artifact_id = ?", cartID, artifactID).First(&entry).Error
	return entry, err
}

func (r *Repository) CreateSiteEntry(entry *models.SiteEntry) error {
	return r.DB.Create(entry).Error
}

func (r *Repository) UpdateSiteEntry(entry *models.SiteEntry) error {
	return r.DB.Save(entry).Error
}

// --- END NEW/SIMPLIFIED CART/ENTRY CRUD METHODS ---

func (r *Repository) DeleteSiteCartSQL(cartID uint, newStatus string) error {
	query := `UPDATE site_carts SET status = $1 WHERE id = $2`

	result, err := r.SQLDB.Exec(query, newStatus, cartID)
	if err != nil {
		return fmt.Errorf("SQL UPDATE execution error: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *Repository) RemoveArtifactFromCart(cartID uint, artifactID uint) error {
	res := r.DB.Where("cart_id = ? AND artifact_id = ?", cartID, artifactID).Delete(&models.SiteEntry{})

	if res.Error != nil {
		return fmt.Errorf("error deleting cart entry: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *Repository) UpdateCartDetails(cartID uint, siteName, comment string) error {
	updates := map[string]interface{}{
		"site_name": siteName,
	}

	res := r.DB.Model(&models.SiteCart{}).Where("id = ?", cartID).Updates(updates)

	if res.Error != nil {
		return fmt.Errorf("error updating cart details: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *Repository) UpdateArtifactQuantityInCart(cartID uint, artifactID uint, quantity int) error {
	res := r.DB.Model(&models.SiteEntry{}).
		Where("cart_id = ? AND artifact_id = ?", cartID, artifactID).
		Update("artifact_quantity", quantity)

	if res.Error != nil {
		return fmt.Errorf("error updating artifact quantity: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}