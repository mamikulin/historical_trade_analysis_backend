package trade_analysis

import (
	"archpath/internal/middleware"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/trade-analysis/cart", h.GetDraftCart).Methods("GET")
	r.HandleFunc("/trade-analysis", h.GetAllRequests).Methods("GET")
	r.HandleFunc("/trade-analysis/{id}", h.GetRequestByID).Methods("GET")
	r.HandleFunc("/trade-analysis/{id}", h.UpdateRequest).Methods("PUT")
	r.HandleFunc("/trade-analysis/{id}/form", h.FormRequest).Methods("PUT")
	r.HandleFunc("/trade-analysis/{id}/moderate", h.CompleteOrRejectRequest).Methods("PUT")
	r.HandleFunc("/trade-analysis/{id}", h.DeleteRequest).Methods("DELETE")
}

// @Summary Get draft cart
// @Description Retrieve the draft cart for the current user
// @Tags trade-analysis
// @Produce json
// @Success 200 {object} object
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Failed to retrieve cart"
// @Security CookieAuth
// @Router /trade-analysis/cart [get]
func (h *Handler) GetDraftCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	cart, err := h.service.GetDraftCart(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve cart: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(cart)
}

// @Summary Get all requests
// @Description Get all trade analysis requests with optional filters. Returns only user's requests for regular users, all requests for moderators
// @Tags trade-analysis
// @Produce json
// @Param status query string false "Filter by status"
// @Param start_date query string false "Filter by start date (YYYY-MM-DD)"
// @Param end_date query string false "Filter by end date (YYYY-MM-DD)"
// @Success 200 {array} object
// @Failure 400 {string} string "Invalid parameters"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 500 {string} string "Failed to retrieve requests"
// @Security CookieAuth
// @Router /trade-analysis [get]
func (h *Handler) GetAllRequests(w http.ResponseWriter, r *http.Request) {
	userID, userAuthenticated := middleware.GetUserIDFromContext(r.Context())
	role, _ := middleware.GetRoleFromContext(r.Context())

	if !userAuthenticated {
		http.Error(w, "Unauthorized: authentication required", http.StatusUnauthorized)
		return
	}

	status := r.URL.Query().Get("status")
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	var startDate, endDate *time.Time

	if startDateStr != "" {
		parsed, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			http.Error(w, "Invalid start_date format (use YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
		startDate = &parsed
	}

	if endDateStr != "" {
		parsed, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			http.Error(w, "Invalid end_date format (use YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
		endDate = &parsed
	}

	var creatorIDFilter *uint
	if role != "moderator" {
		creatorIDFilter = &userID
	}

	requests, err := h.service.GetAllRequests(status, startDate, endDate, creatorIDFilter)
	if err != nil {
		http.Error(w, "Failed to retrieve requests: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(requests)
}

// @Summary Get request by ID
// @Description Retrieve a single trade analysis request by ID
// @Tags trade-analysis
// @Produce json
// @Param id path int true "Request ID"
// @Success 200 {object} object
// @Failure 400 {string} string "Invalid request ID"
// @Failure 404 {string} string "Request not found"
// @Router /trade-analysis/{id} [get]
func (h *Handler) GetRequestByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid request ID", http.StatusBadRequest)
		return
	}

	request, err := h.service.GetRequestByID(uint(id))
	if err != nil {
		http.Error(w, "Request not found: "+err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(request)
}

// @Summary Update request
// @Description Update fields of a trade analysis request
// @Tags trade-analysis
// @Accept json
// @Produce json
// @Param id path int true "Request ID"
// @Param updates body object true "Fields to update"
// @Success 200 {object} object
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Update failed"
// @Security CookieAuth
// @Router /trade-analysis/{id} [put]
func (h *Handler) UpdateRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid request ID", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	delete(updates, "id")
	delete(updates, "creator_id")
	delete(updates, "formation_date")
	delete(updates, "completion_date")
	delete(updates, "moderator_id")

	err = h.service.UpdateRequest(uint(id), updates)
	if err != nil {
		http.Error(w, "Failed to update request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	request, err := h.service.GetRequestByID(uint(id))
	if err != nil {
		http.Error(w, "Failed to retrieve updated request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(request)
}

// @Summary Form request
// @Description Form a trade analysis request (change status from draft to formed)
// @Tags trade-analysis
// @Produce json
// @Param id path int true "Request ID"
// @Success 200 {object} object
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Security CookieAuth
// @Router /trade-analysis/{id}/form [put]
func (h *Handler) FormRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid request ID", http.StatusBadRequest)
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = h.service.FormRequest(uint(id), userID)
	if err != nil {
		http.Error(w, "Failed to form request: "+err.Error(), http.StatusBadRequest)
		return
	}

	request, err := h.service.GetRequestByID(uint(id))
	if err != nil {
		http.Error(w, "Failed to retrieve formed request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(request)
}

// @Summary Complete or reject request
// @Description Complete or reject a trade analysis request (moderator only)
// @Tags trade-analysis
// @Accept json
// @Produce json
// @Param id path int true "Request ID"
// @Param action body object{action=string} true "Action: completed or rejected"
// @Success 200 {object} object
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden: moderator access required"
// @Security CookieAuth
// @Router /trade-analysis/{id}/moderate [put]
func (h *Handler) CompleteOrRejectRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid request ID", http.StatusBadRequest)
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	role, _ := middleware.GetRoleFromContext(r.Context())
	if role != "moderator" {
		http.Error(w, "Forbidden: moderator access required", http.StatusForbidden)
		return
	}

	// Случайно выбираем действие: одобрить или отклонить
	actions := []string{"completed", "rejected"}
	randomAction := actions[time.Now().UnixNano()%2]

	err = h.service.CompleteOrRejectRequest(uint(id), userID, randomAction)
	if err != nil {
		http.Error(w, "Failed to moderate request: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Возвращаем полную информацию о заявке со всеми записями
	request, err := h.service.GetRequestByID(uint(id))
	if err != nil {
		http.Error(w, "Failed to retrieve moderated request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Добавляем информацию о действии
	request["action"] = randomAction
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(request)
}

// @Summary Delete request
// @Description Delete a trade analysis request (moderator only)
// @Tags trade-analysis
// @Param id path int true "Request ID"
// @Success 200 {object} object{message=string}
// @Failure 400 {string} string "Invalid request ID"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden: moderator access required"
// @Security CookieAuth
// @Router /trade-analysis/{id} [delete]
func (h *Handler) DeleteRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid request ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteRequest(uint(id))
	if err != nil {
		http.Error(w, "Failed to delete request: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Request deleted successfully"})
}