package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/benjaminmishra/abios-apis/internal/api"
	"github.com/benjaminmishra/abios-apis/internal/config"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	apiServer := api.New(ctx, cfg)

	go func() {
		if err := apiServer.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %s", err)
		}
	}()

	// wait for context cancellation or signal
	select {
	case <-quit:
		log.Println("Shutdown signal received, initiating graceful shutdown...")
		cancel()
	case <-ctx.Done():
		log.Println("Context done, initiating graceful shutdown...")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := apiServer.Stop(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %s", err)
	}
}
