package application

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"

	"dramastudio/internal/ai/domain"
)

// catalogFile mirrors configs/model-catalog.json — the shipped preset
// matrix as data. Operators edit the file or the registry via the API;
// the code never hardcodes model names.
type catalogFile struct {
	Providers []catalogProvider `json:"providers"`
	Models    []catalogModel    `json:"models"`
	Policies  []catalogPolicy   `json:"policies"`
}

type catalogProvider struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	BaseURL   string `json:"base_url"`
	APIKeyEnv string `json:"api_key_env"`
}

type catalogModel struct {
	ID           string            `json:"id"`
	ProviderID   string            `json:"provider_id"`
	Name         string            `json:"name"`
	Identifier   string            `json:"identifier"`
	Capabilities []string          `json:"capabilities"`
	Modalities   domain.Modalities `json:"modalities"`
	Version      string            `json:"version"`
	Status       string            `json:"status"`
	Pricing      domain.Pricing    `json:"pricing"`
}

type catalogPolicy struct {
	Capability      string               `json:"capability"`
	ProviderID      string               `json:"provider_id"`
	ModelID         string               `json:"model_id"`
	FallbackModels  []domain.ModelTarget `json:"fallback_models"`
	RoutingStrategy string               `json:"routing_strategy"`
}

// SeedCatalog loads the model catalog into an empty registry. It runs
// once: if any provider already exists the seed is skipped so operator
// edits are never overwritten on restart.
func SeedCatalog(ctx context.Context, repo domain.ModelRegistryRepository, path string) (int, error) {
	existing, err := repo.ListProviders(ctx)
	if err != nil {
		return 0, err
	}
	if len(existing) > 0 {
		return 0, nil // registry already populated or customized
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("model catalog: %w", err)
	}
	var cat catalogFile
	if err := json.Unmarshal(data, &cat); err != nil {
		return 0, fmt.Errorf("model catalog parse: %w", err)
	}

	now := time.Now().UTC()
	for _, cp := range cat.Providers {
		if err := repo.SaveProvider(ctx, &domain.Provider{
			ID:           domain.ProviderID(cp.ID),
			Name:         cp.Name,
			Type:         cp.Type,
			BaseURL:      cp.BaseURL,
			APIKeyEnv:    cp.APIKeyEnv,
			HealthStatus: domain.HealthUnknown,
			CreatedAt:    now,
		}); err != nil {
			return 0, fmt.Errorf("seed provider %s: %w", cp.ID, err)
		}
	}

	seeded := 0
	for _, cm := range cat.Models {
		caps := make([]domain.AICapability, len(cm.Capabilities))
		for i, c := range cm.Capabilities {
			caps[i] = domain.AICapability(c)
		}
		if err := repo.SaveModel(ctx, &domain.Model{
			ID:           domain.ModelID(cm.ID),
			Name:         cm.Name,
			ProviderID:   cm.ProviderID,
			Identifier:   cm.Identifier,
			Capabilities: caps,
			Modalities:   cm.Modalities,
			Pricing:      cm.Pricing,
			Version:      cm.Version,
			Status:       domain.ModelStatus(cm.Status),
			CreatedAt:    now,
		}); err != nil {
			return 0, fmt.Errorf("seed model %s: %w", cm.ID, err)
		}
		seeded++
	}

	for _, cp := range cat.Policies {
		if err := repo.UpsertPolicy(ctx, &domain.PolicyRow{
			ID:              "pol_" + uuid.NewString(),
			Scope:           domain.ScopeSystem,
			Capability:      domain.AICapability(cp.Capability),
			ProviderID:      cp.ProviderID,
			ModelID:         cp.ModelID,
			FallbackModels:  cp.FallbackModels,
			RoutingStrategy: cp.RoutingStrategy,
			CreatedAt:       now,
		}); err != nil {
			return 0, fmt.Errorf("seed policy %s: %w", cp.Capability, err)
		}
	}
	return seeded, nil
}
