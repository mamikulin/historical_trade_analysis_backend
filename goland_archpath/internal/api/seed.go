package api

import (
	"archpath/internal/app/models"
	"archpath/internal/app/repository"
)


// SeedData initializes the database with default users and artifacts.
func SeedData(r *repository.Repository) error {
	// 1. Seed Default User
	user := models.User{Login: "user1"}
	r.DB.FirstOrCreate(&user, models.User{Login: "user1"}, &user)

	// 2. Seed Artifacts using the fixed model structure.
	// NOTE: Old 'Period' data is now in 'Description', and 'Region' data is now in 'ProductionCenter'.
	const MinIOBaseURL = "http://localhost:9000/arcpath/"
	artifacts := []models.Artifact{
		{
			Name: "Амфоры аттические",
			Description: "Период: V–IV вв. до н. э.", // Mapped from old Period
			ProductionCenter: "Аттика", // Mapped from old Region
			ExampleLocation: nil, 
			ImageURL: strPtr(MinIOBaseURL + "first.jpg"),
		},
		{
			Name: "Амфоры коринфские",
			Description: "Период: VI–V вв. до н. э.",
			ProductionCenter: "Коринф",
			ExampleLocation: nil,
			ImageURL: strPtr(MinIOBaseURL + "2.jpg"),
		},
		{
			Name: "Финикийские амфоры",
			Description: "Период: VIII–VI вв. до н. э.",
			ProductionCenter: "Финикия",
			ExampleLocation: nil,
			ImageURL: strPtr(MinIOBaseURL + "3.jpeg"),
		},
		{
			Name: "Римские монеты",
			Description: "Период: I–IV вв. н. э.",
			ProductionCenter: "Италия",
			ExampleLocation: nil,
			ImageURL: strPtr(MinIOBaseURL + "4.jpg"),
		},
		{
			Name: "Бронзовые наконечники стрел скифские",
			Description: "Период: VII–III вв. до н. э.",
			ProductionCenter: "Скифия",
			ExampleLocation: nil,
			ImageURL: strPtr(MinIOBaseURL + "5.jpg"),
		},
		{
			Name: "Фибулы латенской культуры",
			Description: "Период: IV–I вв. до н. э.",
			ProductionCenter: "Центральная Европа",
			ExampleLocation: nil,
			ImageURL: strPtr(MinIOBaseURL + "6.jpg"),
		},
	}

	for _, a := range artifacts {
		// Use FirstOrCreate to prevent duplicating data on subsequent runs
		r.DB.FirstOrCreate(&models.Artifact{}, models.Artifact{Name: a.Name}, &a)
	}

	return nil
}
