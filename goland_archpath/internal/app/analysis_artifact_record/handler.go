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
	r.HandleFunc("/analysis-artifact-records/{request_id}/{artifact_id}", h.UpdateRecord).Methods("PUT")
	r.HandleFunc("/analysis-artifact-records/{request_id}/{artifact_id}", h.DeleteRecord).Methods("DELETE")
}

// @Summary Update analysis artifact record
// @Description Update an analysis artifact record (moderator only)
// @Tags analysis-artifact-records
// @Accept json
// @Produce json
// @Param request_id path int true "Request ID"
// @Param artifact_id path int true "Artifact ID"
// @Param updates body object true "Fields to update"
// @Success 200 {object} AnalysisArtifactRecord
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden: moderator access required"
// @Failure 500 {string} string "Update failed"
// @Security CookieAuth
// @Router /analysis-artifact-records/{request_id}/{artifact_id} [put]
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

	delete(updates, "request_id")
	delete(updates, "artifact_id")
	delete(updates, "percentage") // Percentage is calculated, not manually updated

	err = h.service.UpdateRecord(uint(requestID), uint(artifactID), updates)
	if err != nil {
		http.Error(w, "Failed to update record: "+err.Error(), http.StatusInternalServerError)
		return
	}

	record, err := h.service.GetRecordByCompositeKey(uint(requestID), uint(artifactID))
	if err != nil {
		http.Error(w, "Failed to retrieve updated record: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(record)
}

// @Summary Delete analysis artifact record
// @Description Delete an analysis artifact record (moderator only)
// @Tags analysis-artifact-records
// @Param request_id path int true "Request ID"
// @Param artifact_id path int true "Artifact ID"
// @Success 200 {object} object{message=string}
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden: moderator access required"
// @Failure 500 {string} string "Deletion failed"
// @Security CookieAuth
// @Router /analysis-artifact-records/{request_id}/{artifact_id} [delete]
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