package handler

import (
	"encoding/json"
	"net/http"

	"github.com/marmota-alpina/weather-by-cep/cep-weather-service/internal/config"
	"github.com/marmota-alpina/weather-by-cep/cep-weather-service/internal/service"
)

type WeatherHandler struct {
	svc *service.WeatherService
}

func NewWeatherHandler(cfg *config.Config) *WeatherHandler {
	return &WeatherHandler{
		svc: service.NewWeatherService(cfg),
	}
}

type CepRequest struct {
	Cep string `json:"cep"`
}

type Geolocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorResponse{Error: message})
}

func (h *WeatherHandler) HandleWeather(w http.ResponseWriter, r *http.Request) {
	var req CepRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request")
		return
	}

	result, status, err := h.svc.ProcessCep(r.Context(), req.Cep)
	if err != nil {
		writeJSONError(w, status, "Erro ao processar CEP: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(result)
}

func (h *WeatherHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var req Geolocation
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request")
		return
	}

	result, status, err := h.svc.ProcessGeolocation(r.Context(), req.Latitude, req.Longitude)
	if err != nil {
		writeJSONError(w, status, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(result)
}
