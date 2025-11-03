package api

import (
	"archpath/internal/app/analysis_artifact_record"
	"archpath/internal/app/trade_analysis"
	"archpath/internal/app/artifact"
	"archpath/internal/app/user"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// StartServer sets up the router and starts the HTTP server.
func StartServer(artifactService artifact.Service, userService *user.Service, aarService *analysis_artifact_record.Service, taService *trade_analysis.Service) {
	r := mux.NewRouter()

	// Register artifact routes
	artifactHandler := artifact.NewHandler(artifactService)
	artifactHandler.RegisterRoutes(r)

	// Register user routes
	userHandler := user.NewHandler(userService)
	userHandler.RegisterRoutes(r)

	// Register analysis-artifact-record routes
	aarHandler := analysis_artifact_record.NewHandler(aarService)
	aarHandler.RegisterRoutes(r)

    taHandler := trade_analysis.NewHandler(taService)
    taHandler.RegisterRoutes(r)

	addr := ":8000"
	log.Printf("Server listening on %s\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}