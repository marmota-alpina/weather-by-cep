package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	"github.com/marmota-alpina/weather-by-cep/cep-input-service/internal/config"
	"go.opentelemetry.io/otel"
)

type CepRequest struct {
	Cep string `json:"cep"`
}

func isValidCep(cep string) bool {
	match, _ := regexp.MatchString(`^\d{8}$`, cep)
	return match
}

func ProcessCep(cfg *config.Config, cep string, ctx context.Context) (map[string]interface{}, int, error) {
	tr := otel.Tracer("cep-input-service")
	_, span := tr.Start(ctx, "ProcessCep")
	defer span.End()

	if !isValidCep(cep) {
		return nil, http.StatusUnprocessableEntity, errors.New("invalid zipcode")
	}

	body, _ := json.Marshal(CepRequest{Cep: cep})
	url := cfg.WeatherServiceURL + "/weather"

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, http.StatusBadGateway, errors.New("erro ao consultar serviço de clima")
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, http.StatusBadGateway, errors.New(err.Error())
	}

	// Lógica de retorno baseada no status do serviço downstream
	switch resp.StatusCode {
	case http.StatusOK:
		return result, http.StatusOK, nil
	case http.StatusNotFound:
		return nil, http.StatusNotFound, errors.New("can not find zipcode")
	case http.StatusUnprocessableEntity:
		return nil, http.StatusUnprocessableEntity, errors.New("invalid zipcode")
	default:
		return nil, http.StatusBadGateway, errors.New("erro inesperado ao consultar serviço de clima")
	}
}
