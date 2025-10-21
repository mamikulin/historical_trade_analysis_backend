package main

import (
	"archpath/internal/api"
	"archpath/internal/app/repository"
	"archpath/internal/app/service"
	"log"

	"github.com/sirupsen/logrus"
)

const dsn = "host=localhost user=myuser password=mypassword dbname=mydb port=5432 sslmode=disable TimeZone=Europe/Moscow"

func main() {
	log.Println("Application start")

	// 1. Initialize Repository (DB connection and migration)
	repo, err := repository.NewRepository(dsn)
	if err != nil {
		logrus.Fatalf("Failed to initialize repository and connect to DB: %v. Check DSN and 'db' container status.", err)
	}

	// 2. Seed Data
	if err := api.SeedData(repo); err != nil {
		logrus.Errorf("Error seeding initial data: %v", err)
	}

	// 3. Initialize Services
	// Updated: Using NewAnalysisService instead of NewCartService
	analysisService := service.NewAnalysisService(repo)

	// 4. Start Server
	// Passes the correctly named analysis service
	api.StartServer(repo, analysisService) 
	log.Println("Application terminated")
}
