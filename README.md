# Weather by CEP

Este projeto fornece uma API distribuída em dois micro serviços para consultar a previsão do tempo atual a partir de um CEP brasileiro. Ele integra os serviços **FindCEP** (para localização geográfica via CEP) e **OpenWeatherMap** (para dados meteorológicos), com suporte a rastreamento distribuído usando **OpenTelemetry** e **Zipkin**.

---

## 🧱 Arquitetura

- **cep-input-service**: API pública para receber requisições com o CEP.
- **cep-weather-service**: Serviço interno que integra FindCEP e OpenWeatherMap.
- **otel-collector**: Coletor OTEL para rastreamento distribuído (tracing).
- **zipkin**: Interface web para visualização dos traces.

---

## 🚀 Como executar

### Pré-requisitos

- Docker e Docker Compose instalados
- Token de acesso da FindCEP
- API Key do OpenWeatherMap

### Variáveis de Ambiente

No `docker-compose.yml`, você deve substituir:

```yaml
FINDCEP_BASE_URL={SCHEME}://{CLIENT_ID}-{CLIENT_URL_HASH}.api.findcep.com
FINDCEP_REFER={FID}
WEATHER_API_KEY=YOUR WEATHER API KEY
````

### Iniciar os serviços

```bash
docker compose up --build
```

Os serviços estarão disponíveis em:

* `http://localhost:8080` → **cep-input-service**
* `http://localhost:8081` → **cep-weather-service**
* `http://localhost:9411` → **Zipkin**

---
## 📦 Endpoints CEP Input Service
### POST `/cep`
Consulta a previsão do tempo para um CEP.

**Requisição:**

```http
POST /weather HTTP/1.1
Content-Type: application/json

{
  "cep": "01234000"
}
```

**Resposta:**

```json
{
  "city": "São Paulo",
  "temp_C": 18.4,
  "temp_F": 65.2,
  "temp_K": 291.6
}
```
## 📦 Endpoints CEP Weather Service

### POST `/weather`

Consulta a previsão do tempo para um CEP.

**Requisição:**

```http
POST /weather HTTP/1.1
Content-Type: application/json

{
  "cep": "01234000"
}
```

**Resposta:**

```json
{
  "city": "São Paulo",
  "temp_C": 18.4,
  "temp_F": 65.2,
  "temp_K": 291.6
}
```

---

## 🧪 Observabilidade

Com o **otel-collector** e **Zipkin** integrados, você pode visualizar os traces das requisições distribuídas via:

[http://localhost:9411](http://localhost:9411)

---

## 📄 Referências

* [FindCEP API Docs](https://www.findcep.com/docs/index.html)
* [OpenWeatherMap One Call API 3.0](https://openweathermap.org/api/one-call-3)
* [OpenTelemetry Collector](https://opentelemetry.io/docs/collector/)
* [Zipkin](https://zipkin.io/)