package api

import (
	"context"
	"log"
	"net/http"

	"github.com/benjaminmishra/abios-apis/internal/abios"
	"github.com/benjaminmishra/abios-apis/internal/config"
	"github.com/benjaminmishra/abios-apis/internal/service"
	"golang.org/x/time/rate"
)

type Server struct {
	httpServer *http.Server
}

func New(ctx context.Context, cfg *config.Config) *Server {

	client := abios.NewClient(cfg.ApiBaseUrl, cfg.Token, 10, 5, 10)
	liveService := service.NewAbiosLiveService(client)
	handler := NewHandler(ctx, liveService)

	// setup rate limit middleware
	limiter := rate.NewLimiter(5, 10)

	// routes
	mux := http.NewServeMux()
	setupRoutes(mux, handler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: rateLimitMiddleware(mux, limiter),
	}

	return &Server{httpServer: srv}
}

func (s *Server) Start() error {
	log.Println("server listening on :8080")
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	log.Println("shutting down server...")
	return s.httpServer.Shutdown(ctx)
}

func setupRoutes(mux *http.ServeMux, handler *handler) {
	mux.HandleFunc("/series/live", handler.GetLiveSeries)
	mux.HandleFunc("/players/live", handler.GetLivePlayers)
	mux.HandleFunc("/teams/live", handler.GetLiveTeams)
}
