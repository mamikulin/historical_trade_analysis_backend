package api

import (
	"archpath/internal/app/handler"
	"archpath/internal/app/models"
	"archpath/internal/app/repository"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const dsn = "host=localhost user=myuser password=mypassword dbname=mydb port=5432 sslmode=disable TimeZone=Europe/Moscow"

func strPtr(s string) *string {
	return &s
}

func StartServer() {
	log.Println("Server start up")

	repo, err := repository.NewRepository(dsn)
	if err != nil {
		logrus.Fatalf("Ошибка инициализации репозитория и подключения к БД: %v. Проверьте DSN и статус контейнера 'db'.", err)
	}

	if err := SeedData(repo); err != nil {
		logrus.Errorf("Ошибка наполнения БД начальными данными: %v", err)
	}

	h := handler.NewHandler(repo)

	r := gin.Default()

	r.LoadHTMLGlob(filepath.Join("templates", "*.html"))
	r.Static("/static", filepath.Join("resources"))

	r.GET("/", h.GetArtifactTypes)
	r.GET("/artifact/:id", h.GetArtifactTypeDetails)

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
