package server

import (
	"caldave/internal/config"
	"caldave/internal/handlers"
	"caldave/internal/middleware"
	"log"
	"net/http"
)

func Run(cfg *config.Config) error {
	mux := http.NewServeMux()

	mux.Handle("GET /", handlers.HomeHandler())

	loggedMux := middleware.Logging(mux)
	corsLoggedMux := middleware.SetupCORS(loggedMux)

	log.Printf("Starting server on :%s\n", cfg.Port)
	return http.ListenAndServe(":"+cfg.Port, corsLoggedMux)
}
