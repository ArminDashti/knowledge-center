package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ArminDashti/knowledge-center-api/internal/config"
	"github.com/ArminDashti/knowledge-center-api/internal/db"
	"github.com/ArminDashti/knowledge-center-api/internal/handlers"
	"github.com/ArminDashti/knowledge-center-api/internal/scraper"
	"github.com/ArminDashti/knowledge-center-api/internal/store"
	"github.com/ArminDashti/knowledge-center-api/internal/worker"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()

	if err := db.Migrate(ctx, pool); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	st := store.New(pool)
	sc := scraper.New()
	w := worker.New(st, sc, 64)
	w.Start(ctx)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	api := handlers.New(st, w)
	api.Register(r)

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("starting HTTP server on %s", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()
	log.Printf("ready on %s", cfg.HTTPAddr)

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}
