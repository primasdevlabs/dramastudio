package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"dramastudio/configs"
	"dramastudio/internal/appwiring"
	"dramastudio/internal/platform/workflow/temporal"
)

func main() {
	cfg, err := configs.Load()
	if err != nil {
		slog.Error("config", "err", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// With the in-process engine the pipeline runs inside the API process;
	// this binary only exists to host activity workers for a durable engine.
	if cfg.Workflow.Engine != configs.EngineTemporal {
		slog.Info("workflow engine does not need workers", "engine", cfg.Workflow.Engine)
		return
	}

	container, err := appwiring.Build(ctx, cfg)
	if err != nil {
		slog.Error("wiring", "err", err)
		os.Exit(1)
	}
	defer container.Close()

	tc, err := temporal.Dial(ctx, temporal.Config{
		HostPort:  cfg.Workflow.HostPort,
		Namespace: cfg.Workflow.Namespace,
		APIKey:    cfg.Workflow.APIKey,
		UseTLS:    cfg.Workflow.UseTLS,
	})
	if err != nil {
		slog.Error("temporal dial", "err", err)
		os.Exit(1)
	}
	defer tc.Close()

	workers := temporal.NewBuilder(tc, container.Activities).BuildAll()
	for _, w := range workers {
		go func() {
			if err := w.Run(nil); err != nil {
				slog.Error("worker", "err", err)
			}
		}()
	}
	slog.Info("temporal workers running", "queues", len(workers))

	<-ctx.Done()
	for _, w := range workers {
		w.Stop()
	}
	slog.Info("worker stopped")
}
