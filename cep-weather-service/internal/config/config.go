package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	OtelEndpoint   string
	FindCepBaseURL string
	WeatherApiKey  string
	FindCepRefer   string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	port := getEnv("PORT", "8081")
	otelEndpoint := getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://localhost:4317")
	findCepURL := mustGetEnv("FINDCEP_BASE_URL")
	findCepRefer := mustGetEnv("FINDCEP_REFER")
	weatherKey := mustGetEnv("WEATHER_API_KEY")

	return &Config{
		Port:           port,
		OtelEndpoint:   otelEndpoint,
		FindCepBaseURL: findCepURL,
		WeatherApiKey:  weatherKey,
		FindCepRefer:   findCepRefer,
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func mustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Variável de ambiente obrigatória não definida: %s", key)
	}
	return val
}
