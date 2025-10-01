package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	ApiBaseUrl     string
	Token          string
	ReqTimeout     time.Duration
	RateLimitRPS   int
	RateLimitBurst int
}

func LoadConfig() (*Config, error) {
	token := os.Getenv("ABIOS_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("ABIOS_TOKEN not set")
	}

	apiBaseUrl := os.Getenv("ABIOS_API_BASE_URL")
	if apiBaseUrl == "" {
		return nil, fmt.Errorf("ABIOS_API_BASE_URL not set")
	}

	reqTimeoutSecStr := os.Getenv("ABIOS_CLIENT_REQ_TIMEOUT_SEC")
	if reqTimeoutSecStr == "" {
		return nil, fmt.Errorf("ABIOS_CLIENT_REQ_TIMEOUT_SEC not set")
	}
	reqTimeoutSec, err := strconv.Atoi(reqTimeoutSecStr)
	if err != nil {
		return nil, fmt.Errorf("invalid ABIOS_CLIENT_REQ_TIMEOUT_SEC: %v", err)
	}

	rateLimitRPSStr := os.Getenv("ABIOS_CLIENT_RATE_LIMIT_PERSEC")
	if rateLimitRPSStr == "" {
		return nil, fmt.Errorf("ABIOS_CLIENT_RATE_LIMIT_PERSEC not set")
	}
	rateLimitRPS, err := strconv.Atoi(rateLimitRPSStr)
	if err != nil {
		return nil, fmt.Errorf("invalid ABIOS_CLIENT_RATE_LIMIT_PERSEC: %v", err)
	}

	rateLimitBurstStr := os.Getenv("ABIOS_CLIENT_RATE_LIMIT_BURST")
	if rateLimitBurstStr == "" {
		return nil, fmt.Errorf("ABIOS_CLIENT_RATE_LIMIT_BURST not set")
	}
	rateLimitBurst, err := strconv.Atoi(rateLimitBurstStr)
	if err != nil {
		return nil, fmt.Errorf("invalid ABIOS_CLIENT_RATE_LIMIT_BURST: %v", err)
	}

	return &Config{
		ApiBaseUrl:     apiBaseUrl,
		Token:          token,
		ReqTimeout:     time.Duration(reqTimeoutSec) * time.Second,
		RateLimitRPS:   rateLimitRPS,
		RateLimitBurst: rateLimitBurst,
	}, nil
}
