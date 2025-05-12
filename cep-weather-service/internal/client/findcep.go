package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type FindCepClient struct {
	baseURL    string
	refer      string
	httpClient *http.Client
}

func NewFindCepClient(baseURL string, refer string) *FindCepClient {
	return &FindCepClient{
		baseURL: baseURL,
		refer:   refer,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// CepInfo Structs de resposta
type CepInfo struct {
	UF          string `json:"uf"`
	Cidade      string `json:"cidade"`
	Bairro      string `json:"bairro"`
	Logradouro  string `json:"logradouro"`
	CEP         string `json:"cep"`
	Complemento string `json:"complemento"`
}

type Geolocation struct {
	PostalCode string `json:"postal_code"`
	Location   struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"location"`
	Status bool `json:"status"`
}

// CombinedCepData Resposta combinada
type CombinedCepData struct {
	CepInfo
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// GetCombinedCepData Consulta e combina os dados de CEP e localização
func (c *FindCepClient) GetCombinedCepData(cep string) (*CombinedCepData, error) {
	cepURL := fmt.Sprintf("%s/v1/cep/%s.json", c.baseURL, cep)
	geoURL := fmt.Sprintf("%s/v1/geolocation/cep/%s", c.baseURL, cep)

	cepReq, err := http.NewRequest("GET", cepURL, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição para CEP: %w", err)
	}
	cepReq.Header.Set("Referer", c.refer)

	cepResp, err := c.httpClient.Do(cepReq)

	if err != nil {
		return nil, fmt.Errorf("erro na requisição de CEP: %w", err)
	}
	defer cepResp.Body.Close()

	if cepResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro na resposta de CEP: status %d", cepResp.StatusCode)
	}

	var cepInfo CepInfo
	if err := json.NewDecoder(cepResp.Body).Decode(&cepInfo); err != nil {
		return nil, fmt.Errorf("erro ao decodificar CEP: %w", err)
	}

	// Criar e configurar a requisição para geolocalização
	geoReq, err := http.NewRequest("GET", geoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição para geolocalização: %w", err)
	}
	geoReq.Header.Set("Referer", c.refer)

	geoResp, err := c.httpClient.Do(geoReq)
	if err != nil {
		return nil, fmt.Errorf("erro na requisição de geolocalização: %w", err)
	}
	defer geoResp.Body.Close()

	if geoResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro na resposta de geolocalização: status %d", geoResp.StatusCode)
	}

	var geoInfo Geolocation
	if err := json.NewDecoder(geoResp.Body).Decode(&geoInfo); err != nil {
		return nil, fmt.Errorf("falha ao decodificar geolocalização: %w", err)
	}
	if !geoInfo.Status {
		return nil, errors.New("geolocalização não encontrada")
	}

	return &CombinedCepData{
		CepInfo:   cepInfo,
		Latitude:  geoInfo.Location.Lat,
		Longitude: geoInfo.Location.Lon,
	}, nil
}
