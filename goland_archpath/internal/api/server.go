package api

import (
	"archpath/internal/app/analysis_artifact_record"
	"archpath/internal/app/artifact"
	"archpath/internal/app/trade_analysis"
	"archpath/internal/app/user"
	"archpath/internal/middleware"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	
	_ "archpath/docs" // Импорт сгенерированных Swagger документов
)

func StartServer(
	artifactService artifact.Service,
	userService *user.Service,
	aarService *analysis_artifact_record.Service,
	taService *trade_analysis.Service,
) {
	// Получаем JWT конфигурацию из переменных окружения
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Println("WARNING: JWT_SECRET not set, using default (CHANGE IN PRODUCTION!)")
		jwtSecret = "default-secret-key-change-in-production"
	}

	jwtExpiryHours := 24 // Можно также получить из переменной окружения
	jwtExpiry := time.Duration(jwtExpiryHours) * time.Hour

	r := mux.NewRouter()

	// Применяем CORS middleware первым
	r.Use(middleware.CORSMiddleware)

	// Swagger UI (публичный доступ)
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Применяем JWT AuthMiddleware ко всем маршрутам
	r.Use(middleware.AuthMiddleware(jwtSecret))

	api := r.PathPrefix("/api").Subrouter()

	// User routes (публичные и авторизованные)
	userHandler := user.NewHandler(userService, jwtSecret, jwtExpiry)
	// Публичные
	api.HandleFunc("/users/register", userHandler.Register).Methods("POST")
	api.HandleFunc("/users/login", userHandler.Login).Methods("POST")
	// Требуют авторизации
	api.Handle("/users/me", middleware.RequireAuth(http.HandlerFunc(userHandler.GetMe))).Methods("GET")
	api.Handle("/users/me", middleware.RequireAuth(http.HandlerFunc(userHandler.UpdateMe))).Methods("PUT")
	api.Handle("/users/logout", middleware.RequireAuth(http.HandlerFunc(userHandler.Logout))).Methods("POST")

	// Artifact routes
	artifactHandler := artifact.NewHandler(artifactService)
	// Публичные (чтение)
	api.HandleFunc("/artifacts", artifactHandler.GetAll).Methods("GET")
	api.HandleFunc("/artifacts/{id:[0-9]+}", artifactHandler.GetByID).Methods("GET")
	// Требуют авторизации (пользователь)
	api.Handle("/artifacts/{id:[0-9]+}/add-to-analysis", middleware.RequireAuth(http.HandlerFunc(artifactHandler.AddToDraft))).Methods("POST")
	// Только модератор
	api.Handle("/artifacts", middleware.RequireModerator(http.HandlerFunc(artifactHandler.Create))).Methods("POST")
	api.Handle("/artifacts/{id:[0-9]+}", middleware.RequireModerator(http.HandlerFunc(artifactHandler.Update))).Methods("PUT")
	api.Handle("/artifacts/{id:[0-9]+}", middleware.RequireModerator(http.HandlerFunc(artifactHandler.Delete))).Methods("DELETE")
	api.Handle("/artifacts/{id:[0-9]+}/image", middleware.RequireModerator(http.HandlerFunc(artifactHandler.UploadImage))).Methods("POST")

	// Trade Analysis routes
	taHandler := trade_analysis.NewHandler(taService)
	// Публичные (чтение, но с проверкой внутри)
	api.HandleFunc("/trade-analysis/{id}", taHandler.GetRequestByID).Methods("GET")
	// Требуют авторизации (пользователь)
	api.Handle("/trade-analysis/cart", middleware.RequireAuth(http.HandlerFunc(taHandler.GetDraftCart))).Methods("GET")
	api.Handle("/trade-analysis", middleware.RequireAuth(http.HandlerFunc(taHandler.GetAllRequests))).Methods("GET")
	api.Handle("/trade-analysis/{id}", middleware.RequireAuth(http.HandlerFunc(taHandler.UpdateRequest))).Methods("PUT")
	api.Handle("/trade-analysis/{id}/form", middleware.RequireAuth(http.HandlerFunc(taHandler.FormRequest))).Methods("PUT")
	// Только модератор
	api.Handle("/trade-analysis/{id}/moderate", middleware.RequireModerator(http.HandlerFunc(taHandler.CompleteOrRejectRequest))).Methods("PUT")
	api.Handle("/trade-analysis/{id}", middleware.RequireModerator(http.HandlerFunc(taHandler.DeleteRequest))).Methods("DELETE")

	// Analysis Artifact Record routes (только модератор)
	aarHandler := analysis_artifact_record.NewHandler(aarService)
	api.Handle("/analysis-artifact-records/{request_id}/{artifact_id}", middleware.RequireModerator(http.HandlerFunc(aarHandler.UpdateRecord))).Methods("PUT")
	api.Handle("/analysis-artifact-records/{request_id}/{artifact_id}", middleware.RequireModerator(http.HandlerFunc(aarHandler.DeleteRecord))).Methods("DELETE")

	addr := ":8000"
	log.Printf("Server listening on %s\n", addr)
	log.Printf("Swagger UI available at http://localhost%s/swagger/index.html\n", addr)
	
	// Для HTTPS раскомментируйте следующие строки и создайте сертификаты через mkcert
	// certFile := os.Getenv("TLS_CERT")
	// keyFile := os.Getenv("TLS_KEY")
	// if certFile != "" && keyFile != "" {
	// 	log.Printf("Starting HTTPS server on %s\n", addr)
	// 	if err := http.ListenAndServeTLS(addr, certFile, keyFile, r); err != nil {
	// 		log.Fatalf("Server failed: %v", err)
	// 	}
	// } else {
		if err := http.ListenAndServe(addr, r); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	// }
}