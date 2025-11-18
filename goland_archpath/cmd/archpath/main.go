package main

import (
	"archpath/internal/api"
	"archpath/internal/app/analysis_artifact_record"
	"archpath/internal/app/artifact"
	"archpath/internal/app/minio"
	"archpath/internal/app/trade_analysis"
	"archpath/internal/app/user"
	"log"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const dsn = "host=localhost user=myuser password=mypassword dbname=mydb port=5432 sslmode=disable TimeZone=Europe/Moscow"

// @title ArchPath API
// @version 1.0
// @description API для системы управления артефактами и заявками
// @host localhost:8000
// @BasePath /api
// @securityDefinitions.http BearerAuth
// @in header
// @name Authorization
// @scheme bearer
// @bearerFormat JWT
func main() {
	log.Println("Application start")

	// Подключение к БД
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}

	// Миграции
	if err := db.AutoMigrate(
		&user.User{},
		&artifact.Artifact{},
		&analysis_artifact_record.AnalysisArtifactRecord{},
	); err != nil {
		logrus.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Database migrations completed")

	// Инициализация репозиториев и сервисов
	artifactRepo, err := artifact.NewRepository(dsn)
	if err != nil {
		logrus.Fatalf("Failed to initialize artifact repository: %v", err)
	}

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

	artifactService := artifact.NewService(artifactRepo, minioClient)

	userRepo, err := user.NewRepository(dsn)
	if err != nil {
		logrus.Fatalf("Failed to initialize user repository: %v", err)
	}
	userService := user.NewService(userRepo)

	aarRepo := analysis_artifact_record.NewRepository(db)
	aarService := analysis_artifact_record.NewService(aarRepo)

	taRepo := trade_analysis.NewRepository(db)
	taService := trade_analysis.NewService(taRepo)

	log.Println("All services initialized successfully")

	// Запуск сервера (без sessionManager)
	api.StartServer(artifactService, userService, aarService, taService)

	log.Println("Application terminated")
}