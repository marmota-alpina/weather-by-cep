package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/marmota-alpina/weather-by-cep/cep-weather-service/internal/config"
	"github.com/marmota-alpina/weather-by-cep/cep-weather-service/internal/handler"
)

func NewRouter(cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	h := handler.NewWeatherHandler(cfg)
	r.Post("/weather", h.HandleWeather)
	r.Post("/temp", h.Handle)

	return r
}
