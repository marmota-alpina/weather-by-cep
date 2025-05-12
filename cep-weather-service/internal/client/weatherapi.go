package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type WeatherAPIClient struct {
	APIKey string
}

type WeatherResponse struct {
	Current struct {
		Temp float64 `json:"temp"` // temperatura em Celsius (se units=metric)
	} `json:"current"`
}

func NewWeatherAPIClient(apiKey string) *WeatherAPIClient {
	return &WeatherAPIClient{APIKey: apiKey}
}

func (w *WeatherAPIClient) GetCurrentWeather(ctx context.Context, lat, lon float64) (float64, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	url := fmt.Sprintf(
		"https://api.openweathermap.org/data/3.0/onecall?lat=%f&lon=%f&exclude=minutely,hourly,daily,alerts&units=metric&appid=%s",
		lat, lon, w.APIKey,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, errors.New("failed to fetch weather data")
	}

	var weatherResp WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResp); err != nil {
		return 0, err
	}

	return weatherResp.Current.Temp, nil
}
