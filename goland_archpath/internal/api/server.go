package api

import (
	"archpath/internal/app/handler"
	"archpath/internal/app/repository"
	"archpath/internal/app/service"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// strPtr is a utility function to get a string pointer
func strPtr(s string) *string {
	return &s
}

// StartServer sets up the router and starts the HTTP server.
func StartServer(repo *repository.Repository, analysisService *service.AnalysisService) {
	log.Println("Server start up")
	h := handler.NewHandler(repo, analysisService)
	r := gin.Default()

	// Load templates and static files relative to execution path
	r.LoadHTMLGlob(filepath.Join("templates", "*.html"))
	r.Static("/static", filepath.Join("resources"))

	// Artifact Routes
	r.GET("/", h.GetArtifactTypes)
	r.GET("/artifact/:id", h.GetArtifactTypeDetails)

	// Trade Analysis Routes
	// Note: The service/handler logic uses the new 'analysis' naming.
	r.GET("/analysis", h.GetTradeAnalysis)
	r.GET("/analysis/:id", h.GetTradeAnalysis)
	r.POST("/analysis/:id", h.UpdateTradeAnalysis)
	r.POST("/analysis/add", h.AddArtifactToAnalysis)
	r.POST("/analysis/update_quantity", h.UpdateArtifactQuantityInAnalysis)
	r.POST("/analysis/delete", h.DeleteTradeAnalysis)
	r.POST("/analysis/remove", h.RemoveArtifactFromAnalysis)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Server listening on :%s", port)
	r.Run(":" + port)
	log.Println("Server down")
}
