// Command server runs the calculator API and, when STATIC_DIR is set, serves
// the frontend build from the same origin.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sezzle-calculator/internal/httpapi"
)

type config struct {
	port      string
	staticDir string
}

func loadConfig() config {
	cfg := config{port: "8080", staticDir: os.Getenv("STATIC_DIR")}
	if port := os.Getenv("PORT"); port != "" {
		cfg.port = port
	}
	return cfg
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx, loadConfig(), logger); err != nil {
		logger.Error("server stopped", "error", err)
		os.Exit(1)
	}
}

// run serves until ctx is cancelled, then shuts down gracefully.
func run(ctx context.Context, cfg config, logger *slog.Logger) error {
	if cfg.staticDir != "" {
		if _, err := os.Stat(cfg.staticDir); err != nil {
			return fmt.Errorf("STATIC_DIR: %w", err)
		}
	}

	srv := &http.Server{
		Addr:              ":" + cfg.port,
		Handler:           httpapi.NewRouter(logger, cfg.staticDir),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serveErr := make(chan error, 1)
	go func() {
		logger.Info("listening", "addr", srv.Addr, "static_dir", cfg.staticDir)
		serveErr <- srv.ListenAndServe()
	}()

	select {
	case err := <-serveErr:
		return err
	case <-ctx.Done():
	}

	logger.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
