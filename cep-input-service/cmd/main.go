package main

import (
	"log"
	"net/http"

	"github.com/marmota-alpina/weather-by-cep/cep-input-service/internal/config"
	"github.com/marmota-alpina/weather-by-cep/cep-input-service/internal/router"
	"github.com/marmota-alpina/weather-by-cep/cep-input-service/internal/tracing"
)

func main() {
	cfg := config.LoadConfig()

	tp := tracing.InitTracer(cfg) // Passa a config para o tracer
	defer func() { _ = tp.Shutdown(nil) }()

	r := router.NewRouter(cfg) // Também bom passar cfg para o router futuramente
	log.Printf("cep-input-service rodando na porta %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
