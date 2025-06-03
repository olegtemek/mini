package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/olegtemek/mini/internal/config"
	"github.com/olegtemek/mini/internal/repository"
	"github.com/olegtemek/mini/internal/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.New()
	repo := repository.New(ctx, cfg)
	app := server.New(repo)

	serverErr := make(chan error, 1)
	go func() {
		slog.Info("Starting server...")
		if err := app.Start(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		slog.Info("Shutdown signal received")
	case err := <-serverErr:
		slog.Error("Server crashed", "error", err)
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if err := repo.Close(); err != nil {
		slog.Error("Repository close error", "error", err)
	}

	if err := app.Stop(shutdownCtx); err != nil {
		slog.Error("Server shutdown error", "error", err)
	}

	cancel()
	slog.Info("Server stopped gracefully")
}
