package artifact

import (
	"archpath/internal/middleware"
	"encoding/json"
	"log"
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

// @Summary Get all artifacts
// @Description Get all artifacts with optional filters
// @Tags artifacts
// @Produce json
// @Param production_center query string false "Filter by production center"
// @Success 200 {array} Artifact
// @Failure 500 {string} string "Failed to retrieve artifacts"
// @Router /artifacts [get]
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetAll called: %s %s", r.Method, r.URL.Path)
	filters := map[string]interface{}{}

	if center := r.URL.Query().Get("production_center"); center != "" {
		filters["production_center"] = center
	}
	// Support both 'query' and 'search' parameters
	if query := r.URL.Query().Get("query"); query != "" && len(query) > 0 {
		filters["name"] = query
	} else if search := r.URL.Query().Get("search"); search != "" && len(search) > 0 {
		filters["name"] = search
	}

	artifacts, err := h.service.GetAll(filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artifacts)
}

// @Summary Get artifact by ID
// @Description Get a single artifact by ID
// @Tags artifacts
// @Produce json
// @Param id path int true "Artifact ID"
// @Success 200 {object} Artifact
// @Failure 404 {string} string "Artifact not found"
// @Router /artifacts/{id} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetByID called: %s %s, id=%s", r.Method, r.URL.Path, mux.Vars(r)["id"])
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	artifact, err := h.service.GetByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(artifact)
}

// @Summary Create artifact
// @Description Create a new artifact (moderator only)
// @Tags artifacts
// @Accept json
// @Produce json
// @Param artifact body Artifact true "Artifact data"
// @Success 201 {object} Artifact
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden: moderator access required"
// @Failure 500 {string} string "Creation failed"
// @Security CookieAuth
// @Router /artifacts [post]
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
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(artifact)
}

// @Summary Update artifact
// @Description Update an existing artifact (moderator only)
// @Tags artifacts
// @Accept json
// @Produce json
// @Param id path int true "Artifact ID"
// @Param artifact body Artifact true "Artifact data"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden: moderator access required"
// @Failure 500 {string} string "Update failed"
// @Security CookieAuth
// @Router /artifacts/{id} [put]
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

// @Summary Delete artifact
// @Description Delete an artifact (moderator only)
// @Tags artifacts
// @Param id path int true "Artifact ID"
// @Success 204 {string} string "No Content"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden: moderator access required"
// @Failure 500 {string} string "Deletion failed"
// @Security CookieAuth
// @Router /artifacts/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if err := h.service.Delete(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// @Summary Upload artifact image
// @Description Upload an image for an artifact (moderator only)
// @Tags artifacts
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Artifact ID"
// @Param image formData file true "Image file"
// @Success 200 {object} object{image_url=string}
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden: moderator access required"
// @Failure 500 {string} string "Upload failed"
// @Security CookieAuth
// @Router /artifacts/{id}/image [post]
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

// @Summary Add artifact to draft request
// @Description Add an artifact to the user's draft trade analysis request
// @Tags artifacts
// @Accept json
// @Produce json
// @Param id path int true "Artifact ID"
// @Param body body object{quantity=int,comment=string} false "Quantity and comment"
// @Success 201 {object} object
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Failed to add to draft"
// @Security CookieAuth
// @Router /artifacts/{id}/add-to-analysis [post]
func (h *Handler) AddToDraft(w http.ResponseWriter, r *http.Request) {
	log.Printf("AddToDraft called: %s %s", r.Method, r.URL.Path)
	artifactID, _ := strconv.Atoi(mux.Vars(r)["id"])
	log.Printf("Attempting to add artifact ID: %d", artifactID)
	
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	var body struct {
		Quantity int `json:"quantity"`
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
	
	result, err := h.service.AddToDraft(uint(artifactID), userID, body.Quantity)
	if err != nil {
		log.Printf("Failed to add artifact %d to draft: %v", artifactID, err)
		http.Error(w, "Failed to add to draft: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}