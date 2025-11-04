package artifact

import (
	"archpath/internal/app/auth"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/artifacts", h.GetAll).Methods("GET")
	r.HandleFunc("/artifacts/{id:[0-9]+}", h.GetByID).Methods("GET")
	r.HandleFunc("/artifacts", h.Create).Methods("POST")
	r.HandleFunc("/artifacts/{id:[0-9]+}", h.Update).Methods("PUT")
	r.HandleFunc("/artifacts/{id:[0-9]+}", h.Delete).Methods("DELETE")
	r.HandleFunc("/artifacts/{id:[0-9]+}/image", h.UploadImage).Methods("POST")
	r.HandleFunc("/artifacts/{id:[0-9]+}/add-to-analysis", h.AddToDraft).Methods("POST")
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	filters := map[string]interface{}{}

	if isActive := r.URL.Query().Get("is_active"); isActive != "" {
		filters["is_active"] = (isActive == "true")
	}
	if center := r.URL.Query().Get("production_center"); center != "" {
		filters["production_center"] = center
	}

	artifacts, err := h.service.GetAll(filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artifacts)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	artifact, err := h.service.GetByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(artifact)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var artifact Artifact
	if err := json.NewDecoder(r.Body).Decode(&artifact); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if err := h.service.Create(&artifact); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(artifact)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var data Artifact
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if err := h.service.Update(uint(id), data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if err := h.service.Delete(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UploadImage(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "missing image file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	url, err := h.service.UploadImage(uint(id), file, fileHeader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"image_url": url})
}

// AddToDraft adds an artifact to the user's draft request (POST /artifacts/{id}/add-to-analysis)
// Заявка создается пустой с автоматическим указанием создателя, даты создания и статуса
func (h *Handler) AddToDraft(w http.ResponseWriter, r *http.Request) {
	artifactID, _ := strconv.Atoi(mux.Vars(r)["id"])
	
	// Получаем ID текущего пользователя через singleton функцию
	creatorID := auth.GetCurrentUserID()
	
	var body struct {
		Quantity int    `json:"quantity"`
		Comment  string `json:"comment"`
	}
	
	// Default quantity to 1 if not provided
	body.Quantity = 1
	
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		// If body is empty or invalid, just use defaults
		body.Quantity = 1
	}
	
	if body.Quantity <= 0 {
		body.Quantity = 1
	}
	
	result, err := h.service.AddToDraft(uint(artifactID), creatorID, body.Quantity, body.Comment)
	if err != nil {
		http.Error(w, "Failed to add to draft: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}