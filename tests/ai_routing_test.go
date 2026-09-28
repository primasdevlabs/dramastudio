package tests

import (
	"context"
	"testing"

	aiApp "dramastudio/internal/ai/application"
	aiDomain "dramastudio/internal/ai/domain"
	aiInfra "dramastudio/internal/ai/infrastructure/registry"
	"dramastudio/internal/platform/ai/capability"
	"dramastudio/internal/platform/ai/provider"
	aimock "dramastudio/internal/platform/ai/provider/mock"
	"dramastudio/internal/platform/ai/routing"
)

// failFirst is a provider that always returns a retryable error, to prove
// the candidate chain advances.
type failFirst struct{ name string }

func (f failFirst) ProviderName() string { return f.name }
func (f failFirst) Generate(ctx context.Context, model string, spec capability.GenerationSpec) (*capability.Result, error) {
	return nil, &capability.ProviderError{Kind: capability.ErrUnavailable, Provider: f.name, Message: "down", Retryable: true}
}

func seedModels(t *testing.T, repo *aiInfra.InMemoryModelRegistry) {
	ctx := context.Background()
	for _, p := range []*aiDomain.Provider{
		{ID: "bad", Name: "bad", Type: "mock"},
		{ID: "good", Name: "good", Type: "mock"},
	} {
		if err := repo.SaveProvider(ctx, p); err != nil {
			t.Fatalf("provider: %v", err)
		}
	}
	for _, m := range []*aiDomain.Model{
		{ID: "bad.m1", ProviderID: "bad", Identifier: "bad-m1", Status: aiDomain.ModelActive},
		{ID: "good.m1", ProviderID: "good", Identifier: "good-m1", Status: aiDomain.ModelActive},
	} {
		if err := repo.SaveModel(ctx, m); err != nil {
			t.Fatalf("model: %v", err)
		}
	}
}

func TestPolicyFallbackExecutesNextCandidate(t *testing.T) {
	ctx := context.Background()
	repo := aiInfra.NewInMemoryModelRegistry()
	seedModels(t, repo)

	policies := aiApp.NewPolicyService(repo)
	if _, err := policies.SetPolicy(ctx, &aiDomain.PolicyRow{
		Scope:      "system",
		Capability: "image_generation",
		ProviderID: "bad",
		ModelID:    "bad.m1",
		FallbackModels: []aiDomain.ModelTarget{
			{ProviderID: "good", ModelID: "good.m1"},
		},
		RoutingStrategy: aiDomain.StrategyPrimaryFallback,
	}); err != nil {
		t.Fatalf("policy: %v", err)
	}

	r := routing.NewResolver(aiApp.NewRoutingRegistry(repo))
	r.RegisterProvider("bad", provider.Adapter{Sync: failFirst{name: "bad"}})
	r.RegisterProvider("good", provider.Adapter{Sync: aimock.NewSync("good")})

	res, err := r.Resolve(ctx, "image_generation", "system")
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if len(res.Candidates) != 2 {
		t.Fatalf("expected 2 candidates, got %d", len(res.Candidates))
	}

	executor := aiApp.NewExecuteGenerationHandler(repo, r)
	job, err := executor.Handle(ctx, aiApp.ExecuteGenerationCommand{
		ProjectID:  "p1",
		Capability: "image_generation",
		Scope:      "system",
		Prompt:     "test",
	})
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if job.Status != aiDomain.StatusSucceeded {
		t.Fatalf("expected success via fallback, got %s: %s", job.Status, job.Error)
	}
	if job.ProviderID != "good" {
		t.Fatalf("expected fallback provider 'good', got %q", job.ProviderID)
	}
	if job.Attempt != 2 {
		t.Fatalf("expected attempt 2, got %d", job.Attempt)
	}
}

func TestCostOptimizedSortsCandidates(t *testing.T) {
	ctx := context.Background()
	repo := aiInfra.NewInMemoryModelRegistry()
	seedModels(t, repo)
	// make 'good' cheaper than 'bad'
	m, _ := repo.FindModelByID(ctx, "good.m1")
	m.Pricing = aiDomain.Pricing{Currency: "USD", PerUnit: 0.01, Unit: "image"}
	_ = repo.SaveModel(ctx, m)
	b, _ := repo.FindModelByID(ctx, "bad.m1")
	b.Pricing = aiDomain.Pricing{Currency: "USD", PerUnit: 0.10, Unit: "image"}
	_ = repo.SaveModel(ctx, b)

	policies := aiApp.NewPolicyService(repo)
	if _, err := policies.SetPolicy(ctx, &aiDomain.PolicyRow{
		Scope: "system", Capability: "image_generation",
		ProviderID: "bad", ModelID: "bad.m1",
		FallbackModels:  []aiDomain.ModelTarget{{ProviderID: "good", ModelID: "good.m1"}},
		RoutingStrategy: aiDomain.StrategyCostOptimized,
	}); err != nil {
		t.Fatalf("policy: %v", err)
	}
	r := routing.NewResolver(aiApp.NewRoutingRegistry(repo))
	r.RegisterProvider("bad", provider.Adapter{Sync: aimock.NewSync("bad")})
	r.RegisterProvider("good", provider.Adapter{Sync: aimock.NewSync("good")})

	res, err := r.Resolve(ctx, "image_generation", "system")
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if res.Candidates[0].ProviderID != "good" {
		t.Fatalf("cost_optimized should put cheapest first, got %q", res.Candidates[0].ProviderID)
	}
}

func TestPrimaryOnlyBlocksFallback(t *testing.T) {
	ctx := context.Background()
	repo := aiInfra.NewInMemoryModelRegistry()
	seedModels(t, repo)

	policies := aiApp.NewPolicyService(repo)
	if _, err := policies.SetPolicy(ctx, &aiDomain.PolicyRow{
		Scope: "system", Capability: "image_generation",
		ProviderID: "bad", ModelID: "bad.m1",
		FallbackModels:  []aiDomain.ModelTarget{{ProviderID: "good", ModelID: "good.m1"}},
		RoutingStrategy: aiDomain.StrategyPrimaryOnly,
	}); err != nil {
		t.Fatalf("policy: %v", err)
	}
	r := routing.NewResolver(aiApp.NewRoutingRegistry(repo))
	r.RegisterProvider("bad", provider.Adapter{Sync: aimock.NewSync("bad")})
	r.RegisterProvider("good", provider.Adapter{Sync: aimock.NewSync("good")})

	res, err := r.Resolve(ctx, "image_generation", "system")
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if len(res.Candidates) != 1 {
		t.Fatalf("primary_only should yield 1 candidate, got %d", len(res.Candidates))
	}
}

func TestSeedCatalogLoadsMatrix(t *testing.T) {
	ctx := context.Background()
	repo := aiInfra.NewInMemoryModelRegistry()
	n, err := aiApp.SeedCatalog(ctx, repo, "../configs/model-catalog.json")
	if err != nil {
		t.Fatalf("seed catalog: %v", err)
	}
	if n == 0 {
		t.Fatal("expected catalog to seed models")
	}
	// Idempotent: second run must not overwrite.
	n2, err := aiApp.SeedCatalog(ctx, repo, "../configs/model-catalog.json")
	if err != nil {
		t.Fatalf("reseed: %v", err)
	}
	if n2 != 0 {
		t.Fatalf("reseed should be a no-op, seeded %d", n2)
	}
	// Video policy resolves to Veo primary with Seedance/Wan/Kling fallbacks.
	pol, err := repo.ResolvePolicy(ctx, "video_generation", "system")
	if err != nil || pol == nil {
		t.Fatalf("video_generation policy: %v", err)
	}
	if pol.ModelID != "google.veo-3.1" {
		t.Fatalf("expected veo primary, got %s", pol.ModelID)
	}
	if len(pol.FallbackModels) != 3 {
		t.Fatalf("expected 3 fallbacks, got %d", len(pol.FallbackModels))
	}
}
