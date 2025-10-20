package api

import (
	"archpath/internal/app/handler"
	"archpath/internal/app/models"
	"archpath/internal/app/repository"
	"archpath/internal/app/service" // New dependency
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// strPtr is a utility function to get a string pointer
func strPtr(s string) *string {
	return &s
}

// StartServer now accepts repository and cartService as dependencies
func StartServer(repo *repository.Repository, cartService *service.CartService) {
	log.Println("Server start up")

	h := handler.NewHandler(repo, cartService) // Pass both repo and service

	r := gin.Default()

	// Load templates and static files relative to execution path
	r.LoadHTMLGlob(filepath.Join("templates", "*.html"))
	r.Static("/static", filepath.Join("resources"))

	// Artifact Routes
	r.GET("/", h.GetArtifactTypes)
	r.GET("/artifact/:id", h.GetArtifactTypeDetails)

	// Cart Routes
	r.GET("/cart", h.GetSiteCart)
	r.GET("/cart/:id", h.GetSiteCart)
	r.POST("/cart/:id", h.UpdateSiteCart)

	r.POST("/cart/add", h.AddArtifactToCart)
	r.POST("/cart/update_quantity", h.UpdateArtifactQuantityInCart)
	r.POST("/cart/delete", h.DeleteSiteCart)
	r.POST("/cart/remove", h.RemoveArtifactFromCart)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Printf("Server listening on :%s", port)
	r.Run(":" + port)
	log.Println("Server down")
}

// SeedData is now simple and called from main.go
func SeedData(r *repository.Repository) error {
	user := models.User{Login: "user1"}
	r.DB.FirstOrCreate(&user, models.User{Login: "user1"}, &user)

	const MinIOBaseURL = "http://localhost:9000/arcpath/"

	artifacts := []models.Artifact{
		{Name: "Амфоры аттические", Period: "V–IV вв. до н. э.", Region: "Аттика", ImageURL: strPtr(MinIOBaseURL + "first.jpg")},
		{Name: "Амфоры коринфские", Period: "VI–V вв. до н. э.", Region: "Коринф", ImageURL: strPtr(MinIOBaseURL + "2.jpg")},
		{Name: "Финикийские амфоры", Period: "VIII–VI вв. до н. э.", Region: "Финикия", ImageURL: strPtr(MinIOBaseURL + "3.jpeg")},
		{Name: "Римские монеты", Period: "I–IV вв. н. э.", Region: "Италия", ImageURL: strPtr(MinIOBaseURL + "4.jpg")},
		{Name: "Бронзовые наконечники стрел скифские", Period: "VII–III вв. до н. э.", Region: "Скифия", ImageURL: strPtr(MinIOBaseURL + "5.jpg")},
		{Name: "Фибулы латенской культуры", Period: "IV–I вв. до н. э.", Region: "Центральная Европа", ImageURL: strPtr(MinIOBaseURL + "6.jpg")},
	}
	for _, a := range artifacts {
		r.DB.FirstOrCreate(&models.Artifact{}, models.Artifact{Name: a.Name}, &a)
	}

	return nil
}