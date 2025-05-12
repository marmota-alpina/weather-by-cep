package main

import (
	"log"
	"net/http"

	"github.com/marmota-alpina/weather-by-cep/cep-weather-service/internal/config"
	"github.com/marmota-alpina/weather-by-cep/cep-weather-service/internal/router"
	"github.com/marmota-alpina/weather-by-cep/cep-weather-service/internal/tracing"
)

func main() {
	cfg := config.LoadConfig()

	tp := tracing.InitTracer(cfg)
	defer func() { _ = tp.Shutdown(nil) }()

	r := router.NewRouter(cfg)

	log.Println("cep-weather-service rodando na porta", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
