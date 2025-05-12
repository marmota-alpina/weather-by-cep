package config

import (
	"log"
	"os"
)

type Config struct {
	Port              string
	OtelEndpoint      string
	WeatherServiceURL string
}

func LoadConfig() *Config {
	cfg := &Config{
		Port:              getEnv("PORT", "8080"),
		WeatherServiceURL: getEnv("WEATHER_SERVICE_URL", "http://cep-weather-service:8081/weather"),
		OtelEndpoint:      getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "otel-collector:4317"),
	}
	log.Printf("Config carregada: %+v", cfg)
	return cfg
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
