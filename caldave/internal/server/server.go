package server

import (
	"caldave/internal/config"
	"caldave/internal/handlers"
	"caldave/internal/middleware"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func Run(cfg *config.Config, ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("static"))
	wsHandler := handlers.NewWebSocketHandler()

	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))
	mux.Handle("GET /ws", wsHandler.Handler())
	mux.Handle("GET /", handlers.HomeHandler())

	loggedMux := middleware.Logging(mux)
	corsLoggedMux := middleware.SetupCORS(loggedMux)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: corsLoggedMux,
	}

	go func() {
		log.Printf("Starting server on :%s\n", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error listening and serving: %s\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		<-ctx.Done()
		log.Println("Shutting down the server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("error shutting down http server: %s\n", err)
		}
	}()

	wg.Wait()
	return nil
}
