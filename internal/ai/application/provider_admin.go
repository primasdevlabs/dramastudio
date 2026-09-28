package application

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"dramastudio/internal/ai/domain"
	"dramastudio/internal/platform/ai/provider"
	"dramastudio/internal/platform/ai/routing"
)

// ProviderAdminService manages provider rows and their live adapters:
// create (registers the adapter immediately), test-connection, model
// discovery, delete (unregisters). Self-hosted operators add endpoints
// here without touching application code.
type ProviderAdminService struct {
	repo     domain.ModelRegistryRepository
	resolver *routing.Resolver
}

func NewProviderAdminService(repo domain.ModelRegistryRepository, resolver *routing.Resolver) *ProviderAdminService {
	return &ProviderAdminService{repo: repo, resolver: resolver}
}

// CreateProvider persists the row and registers its adapter so it can
// serve traffic immediately. The API key is stored server-side only.
func (s *ProviderAdminService) CreateProvider(ctx context.Context, p *domain.Provider) (*domain.Provider, error) {
	if p.ID == "" {
		p.ID = domain.ProviderID("prov_" + uuid.NewString())
	}
	if p.HealthStatus == "" {
		p.HealthStatus = domain.HealthUnknown
	}
	p.CreatedAt = time.Now().UTC()
	if err := s.repo.SaveProvider(ctx, p); err != nil {
		return nil, err
	}
	if err := s.RegisterAdapter(p); err != nil {
		return p, fmt.Errorf("provider saved but adapter unavailable: %w", err)
	}
	return p, nil
}

// RegisterAdapter builds and registers the provider's adapter from its
// type. Unknown adapter kinds (pending native integrations) return an
// error but leave the provider row usable as registry data.
func (s *ProviderAdminService) RegisterAdapter(p *domain.Provider) error {
	adapter, err := provider.BuildAdapter(p.Type, provider.Config{
		Name:      string(p.ID),
		BaseURL:   p.BaseURL,
		APIKey:    provider.ResolveAPIKey(p.APIKeyEnv, p.APIKey),
		APIKeyEnv: p.APIKeyEnv,
	})
	if err != nil {
		return err
	}
	s.resolver.RegisterProvider(string(p.ID), adapter)
	return nil
}

// TestProvider runs the adapter's connection check and records health.
func (s *ProviderAdminService) TestProvider(ctx context.Context, id string) error {
	p, err := s.repo.FindProviderByID(ctx, domain.ProviderID(id))
	if err != nil {
		return err
	}
	if _, ok := s.resolver.AdapterFor(string(p.ID)); !ok {
		if err := s.RegisterAdapter(p); err != nil {
			return err
		}
	}
	adapter, _ := s.resolver.AdapterFor(string(p.ID))
	tester, ok := adapter.Sync.(provider.Tester)
	if !ok {
		if t2, ok2 := adapter.Async.(provider.Tester); ok2 {
			tester = t2
		} else {
			return fmt.Errorf("provider adapter %q does not support connection testing", p.Type)
		}
	}
	now := time.Now().UTC()
	p.LastHealthCheckAt = &now
	if err := tester.TestConnection(ctx); err != nil {
		p.HealthStatus = domain.HealthDown
		_ = s.repo.SaveProvider(ctx, p)
		return err
	}
	p.HealthStatus = domain.HealthHealthy
	return s.repo.SaveProvider(ctx, p)
}

// DiscoverModels queries the provider endpoint and upserts its model
// identifiers into the registry (status experimental until curated).
func (s *ProviderAdminService) DiscoverModels(ctx context.Context, id string) ([]domain.Model, error) {
	p, err := s.repo.FindProviderByID(ctx, domain.ProviderID(id))
	if err != nil {
		return nil, err
	}
	if _, ok := s.resolver.AdapterFor(string(p.ID)); !ok {
		if err := s.RegisterAdapter(p); err != nil {
			return nil, err
		}
	}
	adapter, _ := s.resolver.AdapterFor(string(p.ID))
	var disc provider.Discoverer
	if d, ok := adapter.Sync.(provider.Discoverer); ok {
		disc = d
	} else if d, ok := adapter.Async.(provider.Discoverer); ok {
		disc = d
	} else {
		return nil, fmt.Errorf("provider adapter %q does not support model discovery", p.Type)
	}
	found, err := disc.DiscoverModels(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]domain.Model, 0, len(found))
	for _, dm := range found {
		m := &domain.Model{
			ID:         domain.ModelID(string(p.ID) + "." + dm.Identifier),
			ProviderID: string(p.ID),
			Name:       dm.Identifier,
			Identifier: dm.Identifier,
			Status:     domain.ModelExperimental, // operator confirms caps
			Metadata:   dm.Metadata,
			CreatedAt:  time.Now().UTC(),
		}
		for _, c := range dm.Capabilities {
			m.Capabilities = append(m.Capabilities, domain.AICapability(c))
		}
		if err := s.repo.SaveModel(ctx, m); err != nil {
			return nil, err
		}
		out = append(out, *m)
	}
	return out, nil
}

// DeleteProvider removes the provider, its models, and its live adapter.
func (s *ProviderAdminService) DeleteProvider(ctx context.Context, id string) error {
	s.resolver.UnregisterProvider(id)
	return s.repo.DeleteProvider(ctx, domain.ProviderID(id))
}

// AdapterAvailable reports whether a provider type has a registered
// adapter factory (drives the UI's "pending integration" state).
func AdapterAvailable(providerType string) bool {
	for _, k := range provider.KnownKinds() {
		if k == providerType {
			return true
		}
	}
	return false
}
