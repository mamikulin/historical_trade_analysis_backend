package analysis_artifact_record

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	// Update a record (quantity, order, etc.) - без PK м-м
	r.HandleFunc("/analysis-artifact-records/{request_id}/{artifact_id}", h.UpdateRecord).Methods("PUT")
	
	// Delete a record from request - без PK м-м
	r.HandleFunc("/analysis-artifact-records/{request_id}/{artifact_id}", h.DeleteRecord).Methods("DELETE")
}

// UpdateRecord updates a record (quantity, order, comment) - без PK м-м (PUT /analysis-artifact-records/{request_id}/{artifact_id})
func (h *Handler) UpdateRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID, err := strconv.ParseUint(vars["request_id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid request_id", http.StatusBadRequest)
		return
	}
	
	artifactID, err := strconv.ParseUint(vars["artifact_id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid artifact_id", http.StatusBadRequest)
		return
	}
	
	var updates map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Prevent updating the composite key fields
	delete(updates, "request_id")
	delete(updates, "artifact_id")
	
	err = h.service.UpdateRecord(uint(requestID), uint(artifactID), updates)
	if err != nil {
		http.Error(w, "Failed to update record: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return the updated record
	record, err := h.service.GetRecordByCompositeKey(uint(requestID), uint(artifactID))
	if err != nil {
		http.Error(w, "Failed to retrieve updated record: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(record)
}

// DeleteRecord deletes a record from the request - без PK м-м (DELETE /analysis-artifact-records/{request_id}/{artifact_id})
func (h *Handler) DeleteRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID, err := strconv.ParseUint(vars["request_id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid request_id", http.StatusBadRequest)
		return
	}
	
	artifactID, err := strconv.ParseUint(vars["artifact_id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid artifact_id", http.StatusBadRequest)
		return
	}
	
	err = h.service.DeleteRecord(uint(requestID), uint(artifactID))
	if err != nil {
		http.Error(w, "Failed to delete record: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Record deleted successfully"})
}