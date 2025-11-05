package main

import (
	"archpath/internal/api"
	"archpath/internal/app/analysis_artifact_record"
	"archpath/internal/app/artifact"
	"archpath/internal/app/minio"
	"archpath/internal/app/session"
	"archpath/internal/app/trade_analysis"
	"archpath/internal/app/user"
	"log"
	"time"

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
// @securityDefinitions.apikey CookieAuth
// @in cookie
// @name session_id
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

	// Инициализация Session Manager
	sessionManager, err := session.NewManager("localhost:6379", 24*time.Hour)
	if err != nil {
		logrus.Fatalf("Failed to initialize session manager: %v", err)
	}
	log.Println("Session manager initialized")

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

	// Запуск сервера
	api.StartServer(artifactService, userService, aarService, taService, sessionManager)

	log.Println("Application terminated")
}