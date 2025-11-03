package api

import (
	"archpath/internal/app/artifact"
)

func strPtr(s string) *string { return &s }

func SeedData(r *artifact.Repository) error {
	const MinIOBaseURL = "http://localhost:9000/arcpath/"

	artifacts := []artifact.Artifact{
		{
			Name:             "Амфоры аттические",
			Description:      "Период: V–IV вв. до н. э.",
			ProductionCenter: "Аттика",
			ImageURL:         strPtr(MinIOBaseURL + "first.jpg"),
		},
		{
			Name:             "Амфоры коринфские",
			Description:      "Период: VI–V вв. до н. э.",
			ProductionCenter: "Коринф",
			ImageURL:         strPtr(MinIOBaseURL + "2.jpg"),
		},
		{
			Name:             "Финикийские амфоры",
			Description:      "Период: VIII–VI вв. до н. э.",
			ProductionCenter: "Финикия",
			ImageURL:         strPtr(MinIOBaseURL + "3.jpeg"),
		},
		{
			Name:             "Римские монеты",
			Description:      "Период: I–IV вв. н. э.",
			ProductionCenter: "Италия",
			ImageURL:         strPtr(MinIOBaseURL + "4.jpg"),
		},
		{
			Name:             "Бронзовые наконечники стрел скифские",
			Description:      "Период: VII–III вв. до н. э.",
			ProductionCenter: "Скифия",
			ImageURL:         strPtr(MinIOBaseURL + "5.jpg"),
		},
		{
			Name:             "Фибулы латенской культуры",
			Description:      "Период: IV–I вв. до н. э.",
			ProductionCenter: "Центральная Европа",
			ImageURL:         strPtr(MinIOBaseURL + "6.jpg"),
		},
	}

	for _, a := range artifacts {
		r.DB.FirstOrCreate(&a, artifact.Artifact{Name: a.Name})
	}

	return nil
}
