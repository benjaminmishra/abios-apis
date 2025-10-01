package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/benjaminmishra/abios-apis/internal/service"
)

type handler struct {
	rootCtx     context.Context
	liveService service.LiveService
}

func NewHandler(ctx context.Context, s service.LiveService) *handler {
	return &handler{
		rootCtx:     ctx,
		liveService: s,
	}
}

func (h *handler) GetLiveSeries(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := h.liveService.GetLiveSeries(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(data) == 0 {
		http.Error(w, "No live series found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, data)
}

func (h *handler) GetLivePlayers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := h.liveService.GetLivePlayers(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(data) == 0 {
		http.Error(w, "No live players found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, data)
}

func (h *handler) GetLiveTeams(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := h.liveService.GetLiveTeams(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(data) == 0 {
		http.Error(w, "No live teams found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, data)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
