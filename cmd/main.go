// Command go-v1 runs the API server.
//
// Wiring order per the v1 guide:
//  1. load env
//  2. init DB connection
//  3. init repository
//  4. init handler (inject repo)
//  5. mount routes + middleware
//  6. start server
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/garrychrstn/top-go/internal/database"
	"github.com/garrychrstn/top-go/internal/handler"
	"github.com/garrychrstn/top-go/internal/repository"
	"github.com/garrychrstn/top-go/internal/util"
)

type config struct {
	Env         string
	Port        string
	DatabaseURL string
	CORSOrigins []string
	AuthToken   string
}

func loadConfig() config {
	return config{
		Env:         envOr("ENV", "development"),
		Port:        envOr("PORT", "8080"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		CORSOrigins: splitList(envOr("CORS_ALLOWED_ORIGINS", "*")),
		AuthToken:   os.Getenv("AUTH_BEARER_TOKEN"),
	}
}

func main() {
	// 1. Load env.
	if err := util.LoadDotEnv(".env"); err != nil {
		slog.Error("load env", "error", err)
		os.Exit(1)
	}
	cfg := loadConfig()

	if cfg.DatabaseURL == "" {
		slog.Error("DATABASE_URL is not set — copy .env.example to .env and fill it in")
		os.Exit(1)
	}

	// Structured logging: text for dev, JSON for production.
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	var base slog.Handler = slog.NewTextHandler(os.Stderr, opts)
	if cfg.Env == "production" {
		base = slog.NewJSONHandler(os.Stderr, opts)
	}
	slog.SetDefault(slog.New(base))

	// Graceful shutdown on SIGINT/SIGTERM.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 2. Init DB connection + migrations.
	pool, err := database.New(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("connect database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := database.Migrate(ctx, pool); err != nil {
		slog.Error("migrate database", "error", err)
		os.Exit(1)
	}

	// 3. Init repository, 4. init handler (inject repo).
	handlers := []handler.Registrar{
		handler.NewUserHandler(repository.NewUserRepository(pool)),
	}

	// 5. Mount routes + middleware, 6. start server.
	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           handler.NewRouter(handlers, cfg.CORSOrigins),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		slog.Info("server listening", "addr", srv.Addr, "env", cfg.Env)
		serverErr <- srv.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	case <-ctx.Done():
		slog.Info("shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			slog.Error("graceful shutdown failed", "error", err)
			os.Exit(1)
		}
	}
	slog.Info("server stopped")
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func splitList(s string) []string {
	var out []string
	for _, part := range strings.Split(s, ",") {
		if part = strings.TrimSpace(part); part != "" {
			out = append(out, part)
		}
	}
	return out
}
