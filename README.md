# Abios Api Wrapper

This service wraps the Abios public Atlas API with a lightweight HTTP server that exposes curated live esports data for series, players, and teams. It provides a simple façade for clients that only need the live endpoints without handling the full Abios schema or authentication workflow.

## Getting Started

### Depedencipes
- Install Go 1.24.3 or newer.
- Clone this repository and open it in your terminal.
- Run `go mod download` if dependencies are not yet cached.

### Run The Server
- Start locally with `go run ./cmd/server`.
- Override the default configuration by exporting:
  - `ABIOS_TOKEN`
  - `ABIOS_API_BASE_URL`
  - `ABIOS_CLIENT_REQ_TIMEOUT_SEC`
  - `ABIOS_CLIENT_RATE_LIMIT_PERSEC`
  - `ABIOS_CLIENT_RATE_LIMIT_BURST`
- The server listens on `http://localhost:8080` and serves:
  - `GET /series/live`
  - `GET /players/live`
  - `GET /teams/live`

## Run Tests
- Execute the full suite with `go test ./...`.
- Add `-v` for verbose output when investigating failures.

## Rate Limits
- Incoming HTTP traffic is shaped by `golang.org/x/time/rate` with a default of 5 requests per second and a burst of 10.
- The Abios client honors the same RPS and burst values when calling upstream endpoints.
- Adjust limits via `ABIOS_CLIENT_RATE_LIMIT_PERSEC` and `ABIOS_CLIENT_RATE_LIMIT_BURST` to match production quotas.
- Requests exceeding the limiter receive HTTP 429 responses.

## Possible Improvements
- Externalize secrets and environment defaults into a configuration file or secret manager.
- Add structured logging and correlation IDs for tracing upstream calls.
- Expand integration tests that hit a mocked Abios server to validate end-to-end behaviour.
- Introduce caching of roster and team lookups to reduce duplicate upstream requests.

## Docker Run 

```bash

docker build -t "abios-api:latest" .

docker run --rm -it \
  -p 8080:8080 \
  -e ABIOS_TOKEN="<your_abios_token>" \
  abios-api

```

