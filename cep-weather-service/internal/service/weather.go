package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/marmota-alpina/weather-by-cep/cep-weather-service/internal/client"
	"github.com/marmota-alpina/weather-by-cep/cep-weather-service/internal/config"
	"go.opentelemetry.io/otel"
)

// WeatherService coordena as chamadas para FindCep e WeatherAPI
type WeatherService struct {
	findCepClient    *client.FindCepClient
	weatherApiClient *client.WeatherAPIClient
}

// NewWeatherService instancia o serviço com os clients configurados
func NewWeatherService(cfg *config.Config) *WeatherService {
	return &WeatherService{
		findCepClient:    client.NewFindCepClient(cfg.FindCepBaseURL, cfg.FindCepRefer),
		weatherApiClient: client.NewWeatherAPIClient(cfg.WeatherApiKey),
	}
}

// WeatherResult define o JSON de resposta
type WeatherResult struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type CurrentWeather struct {
	Dt         int64     `json:"dt"`
	Sunrise    int64     `json:"sunrise"`
	Sunset     int64     `json:"sunset"`
	Temp       float64   `json:"temp"`
	FeelsLike  float64   `json:"feels_like"`
	Pressure   int       `json:"pressure"`
	Humidity   int       `json:"humidity"`
	DewPoint   float64   `json:"dew_point"`
	UVI        float64   `json:"uvi"`
	Clouds     int       `json:"clouds"`
	Visibility int       `json:"visibility"`
	WindSpeed  float64   `json:"wind_speed"`
	WindDeg    int       `json:"wind_deg"`
	WindGust   float64   `json:"wind_gust"`
	Weather    []Weather `json:"weather"`
}

type Weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// isValidCep valida se o CEP tem exatamente 8 dígitos numéricos
func isValidCep(cep string) bool {
	cep = strings.TrimSpace(cep)
	match, _ := regexp.MatchString(`^\d{8}$`, cep)
	return match
}

// ProcessCep realiza a lógica principal de validação, consulta e agregação de dados
func (s *WeatherService) ProcessCep(ctx context.Context, cep string) (*WeatherResult, int, error) {
	tracer := otel.Tracer("cep-weather-service")
	ctx, span := tracer.Start(ctx, "ProcessCep")
	defer span.End()

	if !isValidCep(cep) {
		return nil, http.StatusUnprocessableEntity, errors.New("invalid zipcode")
	}

	data, err := s.findCepClient.GetCombinedCepData(cep)
	if err != nil {
		return nil, http.StatusNotFound, errors.New(err.Error())
	}

	tempC, err := s.weatherApiClient.GetCurrentWeather(ctx, data.Latitude, data.Longitude)
	if err != nil {
		return nil, http.StatusBadGateway, fmt.Errorf("weather API error: %w", err)
	}

	result := &WeatherResult{
		City:  data.Cidade,
		TempC: round(tempC, 1),
		TempF: round(celsiusToFahrenheit(tempC), 1),
		TempK: round(celsiusToKelvin(tempC), 1),
	}

	return result, http.StatusOK, nil
}

// ProcessGeolocation ProcessCep realiza a lógica principal de validação, consulta e agregação de dados
func (s *WeatherService) ProcessGeolocation(ctx context.Context, lat float64, lon float64) (float64, int, error) {
	tracer := otel.Tracer("cep-weather-service")
	ctx, span := tracer.Start(ctx, "ProcessCep")
	defer span.End()

	tempC, err := s.weatherApiClient.GetCurrentWeather(ctx, lat, lon)
	if err != nil {
		return 0, http.StatusBadGateway, fmt.Errorf("weather API error: %w", err)
	}

	return tempC, http.StatusOK, nil
}

// celsiusToFahrenheit converte °C para °F
func celsiusToFahrenheit(c float64) float64 {
	return (c * 9.0 / 5.0) + 32.0
}

// celsiusToKelvin converte °C para K
func celsiusToKelvin(c float64) float64 {
	return c + 273.15
}

// Round arredonda o float para a precisão desejada
func round(value float64, precision int) float64 {
	format := fmt.Sprintf("%%.%df", precision)
	str := fmt.Sprintf(format, value)
	var out float64
	fmt.Sscanf(str, "%f", &out)
	return out
}
