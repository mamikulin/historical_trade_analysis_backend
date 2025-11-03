package main

import (
	"archpath/internal/api"
	"archpath/internal/app/analysis_artifact_record"
	"archpath/internal/app/trade_analysis"
	"archpath/internal/app/artifact"
	"archpath/internal/app/user"
	"archpath/internal/app/minio"
	"log"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const dsn = "host=localhost user=myuser password=mypassword dbname=mydb port=5432 sslmode=disable TimeZone=Europe/Moscow"

func main() {
	log.Println("Application start")

	// Initialize shared database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := db.AutoMigrate(
		&user.User{},
		&artifact.Artifact{}, // Add your artifact model here
		&analysis_artifact_record.AnalysisArtifactRecord{},
	); err != nil {
		logrus.Fatalf("Failed to run migrations: %v", err)
	}

	// 1. Initialize artifact repository
	artifactRepo, err := artifact.NewRepository(dsn)
	if err != nil {
		logrus.Fatalf("Failed to initialize artifact repository: %v", err)
	}

	// 2. Initialize MinIO
	minioClient, err := minio.NewMinioClient(
		"localhost:9000",
		"minio",
		"minio124",
		"archpath",
		false,
	)
	if err != nil {
		logrus.Fatalf("Failed to initialize MinIO: %v", err)
	}

	// 3. Initialize artifact service
	artifactService := artifact.NewService(artifactRepo, minioClient)

	// 4. Initialize user repository and service
	userRepo, err := user.NewRepository(dsn)
	if err != nil {
		logrus.Fatalf("Failed to initialize user repository: %v", err)
	}
	userService := user.NewService(userRepo)

	// 5. Initialize analysis-artifact-record repository and service
	aarRepo := analysis_artifact_record.NewRepository(db)
	aarService := analysis_artifact_record.NewService(aarRepo)

    taRepo := trade_analysis.NewRepository(db)
    taService := trade_analysis.NewService(taRepo)

	// 6. Start HTTP server
	api.StartServer(artifactService, userService, aarService, taService)
	
	log.Println("Application terminated")
}