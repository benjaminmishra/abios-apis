FROM golang:1.24.3-bookworm AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

    COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o /out/abios-api ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot

# Setting env variables with defaults here for now.
ENV ABIOS_API_BASE_URL="https://atlas.abiosgaming.com/v3" \
    ABIOS_CLIENT_REQ_TIMEOUT_SEC="10" \
    ABIOS_CLIENT_RATE_LIMIT_PERSEC="5" \
    ABIOS_CLIENT_RATE_LIMIT_BURST="10" 

COPY --from=builder /out/abios-api /usr/bin/abios-api
USER nonroot:nonroot

ENTRYPOINT ["/usr/bin/abios-api"]
