package repository

import (
	"archpath/internal/app/models"
	"database/sql"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Repository holds the GORM and raw SQL database connections.
type Repository struct {
	DB    *gorm.DB
	SQLDB *sql.DB
}

// NewRepository connects to the database and performs AutoMigration for all models.
func NewRepository(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("database connection error: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("error getting *sql.DB: %w", err)
	}

	// Updated AutoMigrate calls to use the new model names
	err = db.AutoMigrate(&models.User{}, &models.Artifact{}, &models.TradeAnalysis{}, &models.AnalysisArtifactRecord{})
	if err != nil {
		return nil, fmt.Errorf("migration error: %w", err)
	}

	return &Repository{DB: db, SQLDB: sqlDB}, nil
}
