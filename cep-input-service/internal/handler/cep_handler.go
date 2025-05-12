package handler

import (
	"encoding/json"
	"net/http"

	"github.com/marmota-alpina/weather-by-cep/cep-input-service/internal/config"
	"github.com/marmota-alpina/weather-by-cep/cep-input-service/internal/service"
)

type Handler struct {
	Cfg *config.Config
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{Cfg: cfg}
}

type errorResponse struct {
	Error string `json:"error"`
}

func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorResponse{Error: msg})
}

func (h *Handler) HandleCep(w http.ResponseWriter, r *http.Request) {
	var input service.CepRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSONError(w, http.StatusUnprocessableEntity, "invalid zipcode")
		return
	}

	result, status, err := service.ProcessCep(h.Cfg, input.Cep, r.Context())
	if err != nil {
		writeJSONError(w, status, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
