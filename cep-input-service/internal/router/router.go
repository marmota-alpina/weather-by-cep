package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/marmota-alpina/weather-by-cep/cep-input-service/internal/config"
	"github.com/marmota-alpina/weather-by-cep/cep-input-service/internal/handler"
)

func NewRouter(cfg *config.Config) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/cep", handler.NewHandler(cfg).HandleCep)
	return r
}
