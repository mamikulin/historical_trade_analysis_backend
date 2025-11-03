package trade_analysis

import (
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

func (h *Handler) GetDraftCart(w http.ResponseWriter, r *http.Request) {
	creatorID := uint(1)
	
	cart, err := h.service.GetDraftCart(creatorID)
	if err != nil {
		http.Error(w, "Failed to retrieve cart: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(cart)
}

func (h *Handler) GetAllRequests(w http.ResponseWriter, r *http.Request) {
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
	
	requests, err := h.service.GetAllRequests(status, startDate, endDate)
	if err != nil {
		http.Error(w, "Failed to retrieve requests: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(requests)
}

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

func (h *Handler) FormRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid request ID", http.StatusBadRequest)
		return
	}
	
	creatorID := uint(1)
	
	err = h.service.FormRequest(uint(id), creatorID)
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

func (h *Handler) CompleteOrRejectRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid request ID", http.StatusBadRequest)
		return
	}
	
	var body struct {
		Action string `json:"action"` 
	}
	
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	moderatorID := uint(2)
	
	err = h.service.CompleteOrRejectRequest(uint(id), moderatorID, body.Action)
	if err != nil {
		http.Error(w, "Failed to moderate request: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	request, err := h.service.GetRequestByID(uint(id))
	if err != nil {
		http.Error(w, "Failed to retrieve moderated request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(request)
}

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