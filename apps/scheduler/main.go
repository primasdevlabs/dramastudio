package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"dramastudio/configs"
	"dramastudio/internal/appwiring"
)

// scheduler runs periodic sweeps: due publications, stuck generation jobs.
func main() {
	cfg, err := configs.Load()
	if err != nil {
		slog.Error("config", "err", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	container, err := appwiring.Build(ctx, cfg)
	if err != nil {
		slog.Error("wiring", "err", err)
		os.Exit(1)
	}
	defer container.Close()

	interval := 30 * time.Second
	if v := os.Getenv("SCHEDULER_INTERVAL"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			interval = d
		}
	}

	slog.Info("scheduler running", "interval", interval, "store", cfg.Store)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	sweep := func() {
		published, failed, err := container.Publishing.PublishDue(ctx, time.Now().UTC())
		if err != nil {
			slog.Error("publish sweep", "err", err)
			return
		}
		if published+failed > 0 {
			slog.Info("publish sweep", "published", published, "failed", failed)
		}
	}
	sweep()

	for {
		select {
		case <-ctx.Done():
			slog.Info("scheduler stopped")
			return
		case <-ticker.C:
			sweep()
		}
	}
}
