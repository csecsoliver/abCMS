package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/csecsoliver/abCMS/internal/auth"
	"github.com/csecsoliver/abCMS/internal/config"
	"github.com/csecsoliver/abCMS/internal/db"
	"github.com/csecsoliver/abCMS/internal/handler"
	"github.com/csecsoliver/abCMS/internal/media"
)

func main() {
	cfg := config.Load()

	store, err := db.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer store.Close()

	sessions := auth.NewSessionManager(cfg.SessionSecret, cfg.SessionTTL, cfg.SecureCookies)
	mediaSvc := media.NewService(cfg.UploadDir)
	h := handler.New(cfg, store, sessions, mediaSvc)

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           h.Routes(),
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("abCMS listening on %s (db=%s, uploads=%s)", cfg.Addr, cfg.DBPath, cfg.UploadDir)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down abCMS...")
	h.Close()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown error: %v", err)
	}
}
