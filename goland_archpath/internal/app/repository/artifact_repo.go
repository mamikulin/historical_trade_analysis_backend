package repository

import (
	"archpath/internal/app/models"
	"strings"
)

// GetArtifacts retrieves a list of active artifacts, optionally filtered by name search.
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

// GetArtifactByID retrieves a single artifact by its ID.
func (r *Repository) GetArtifactByID(id uint) (models.Artifact, error) {
	var artifact models.Artifact
	if err := r.DB.First(&artifact, id).Error; err != nil {
		return models.Artifact{}, err
	}
	return artifact, nil
}
