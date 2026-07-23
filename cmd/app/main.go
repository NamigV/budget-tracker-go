package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NamigV/budget-tracker-go/internal/config"
	"github.com/NamigV/budget-tracker-go/internal/database"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DB)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	log.Println("database connected")

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("get sql.DB failed: %v", err)
	}
	defer func() { _ = sqlDB.Close() }()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handleHealth)

	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: mux,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("server starting on %s", cfg.Addr)
		serverErrors <- srv.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("server error: %v", err)

	case sig := <-shutdown:
		log.Printf("received signal %v, starting graceful shutdown", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("graceful shutdown timed out (%v), forcing close", err)
			_ = srv.Close()
		}

		log.Println("server stopped")
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok\n"))
}
