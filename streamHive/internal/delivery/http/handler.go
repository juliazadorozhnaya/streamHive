package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"streamHive/streamHive/internal/domain"
)

type Handler struct {
	uc domain.UseCase
}

func NewHandler(uc domain.UseCase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) StartRecording(w http.ResponseWriter, r *http.Request) {
	var req struct {
		StreamURL string `json:"stream_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	jobID, err := h.uc.StartRecording(r.Context(), req.StreamURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"job_id": jobID})
}

func (h *Handler) StopRecording(w http.ResponseWriter, r *http.Request) {
	var req struct {
		JobID string `json:"job_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if err := h.uc.StopRecording(req.JobID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetJobStatus(w http.ResponseWriter, r *http.Request) {
	jobID := chi.URLParam(r, "jobID")
	status, err := h.uc.GetJobStatus(jobID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": status})
}
