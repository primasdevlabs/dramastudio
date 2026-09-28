// Command seed creates a development organization and owner user so a fresh
// environment has known credentials to log in with. Safe to re-run — it
// exits cleanly when the account already exists.
//
//	go run ./apps/seed
//
// Overrides (all optional):
//
//	SEED_ORG, SEED_NAME, SEED_EMAIL, SEED_PASSWORD
package main

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"dramastudio/configs"
	"dramastudio/internal/appwiring"
	"dramastudio/internal/identity/domain"
)

func main() {
	orgName := envOr("SEED_ORG", "DramaStudio Dev")
	name := envOr("SEED_NAME", "Dev Owner")
	email := envOr("SEED_EMAIL", "dev@dramastudio.local")
	password := envOr("SEED_PASSWORD", "dramastudio-dev")

	cfg, err := configs.Load()
	if err != nil {
		slog.Error("config", "err", err)
		os.Exit(1)
	}

	container, err := appwiring.Build(context.Background(), cfg)
	if err != nil {
		slog.Error("build", "err", err)
		os.Exit(1)
	}
	defer container.Close()

	u, org, err := container.Identity.Register(context.Background(), orgName, email, name, password)
	if errors.Is(err, domain.ErrEmailTaken) {
		slog.Info("seed user already exists — nothing to do", "email", email)
		return
	}
	if err != nil {
		slog.Error("seed", "err", err)
		os.Exit(1)
	}

	slog.Info("seeded dev account",
		"org", org.Name,
		"org_id", org.ID,
		"user", u.Name,
		"email", u.Email,
		"role", string(u.Role),
	)
	slog.Info("login credentials", "email", email, "password", password)
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
