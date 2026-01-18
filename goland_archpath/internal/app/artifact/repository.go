package artifact

import (
    "fmt"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

type Repository struct {
    DB *gorm.DB
}

func NewRepository(dsn string) (*Repository, error) {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("failed to connect to DB: %w", err)
    }

    if err := db.AutoMigrate(&Artifact{}); err != nil {
        return nil, fmt.Errorf("failed to migrate: %w", err)
    }

    return &Repository{DB: db}, nil
}

func (r *Repository) GetAll(filters map[string]interface{}) ([]Artifact, error) {
	var artifacts []Artifact
	query := r.DB.Model(&Artifact{})

	for key, value := range filters {
		if key == "name" || key == "description" {
			// Use LIKE for text search
			query = query.Where(key+" ILIKE ?", "%"+value.(string)+"%")
		} else {
			query = query.Where(key+" = ?", value)
		}
	}

	if err := query.Find(&artifacts).Error; err != nil {
		return nil, err
	}
	return artifacts, nil
}


func (r *Repository) GetByID(id uint) (*Artifact, error) {
    var artifact Artifact
    if err := r.DB.Where("id = ?", id).First(&artifact).Error; err != nil {
        return nil, err
    }
    return &artifact, nil
}

func (r *Repository) Create(a *Artifact) error {
    return r.DB.Create(a).Error
}

func (r *Repository) Update(id uint, data Artifact) error {
    return r.DB.Model(&Artifact{}).Where("id = ?", id).Updates(data).Error
}

func (r *Repository) Delete(id uint) error {
    return r.DB.Delete(&Artifact{}, id).Error
}
