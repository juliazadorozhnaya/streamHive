package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(h *Handler) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/start", h.StartRecording)
	r.Post("/stop", h.StopRecording)
	r.Get("/status/{jobID}", h.GetJobStatus)

	return r
}
