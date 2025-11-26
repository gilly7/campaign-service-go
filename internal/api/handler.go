package api

import (
	"encoding/json"
	"net/http"

	"campaign-service/internal/campaign"
	"campaign-service/internal/models"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *campaign.Service
}

func NewHandler(service *campaign.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Router() http.Handler {
	r := chi.NewRouter()
	r.Post("/campaigns", h.CreateCampaign)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("OK")) })
	return r
}

func (h *Handler) CreateCampaign(w http.ResponseWriter, r *http.Request) {
	var req models.CreateCampaignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	campaign, err := h.service.CreateCampaign(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(campaign)
}
