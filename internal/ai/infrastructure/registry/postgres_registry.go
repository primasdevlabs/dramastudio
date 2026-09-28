package registry

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"

	"dramastudio/internal/ai/domain"
	"dramastudio/internal/platform/database/postgres"
)

type PostgresModelRegistry struct {
	q postgres.Querier
}

func NewPostgresModelRegistry(q postgres.Querier) *PostgresModelRegistry {
	return &PostgresModelRegistry{q: q}
}

func mustJSON(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}

// --- Providers ---

func (r *PostgresModelRegistry) SaveProvider(ctx context.Context, p *domain.Provider) error {
	if p.HealthStatus == "" {
		p.HealthStatus = domain.HealthUnknown
	}
	_, err := r.q.Exec(ctx, `
		INSERT INTO ai.providers (id, name, type, base_url, api_key_env, api_key, health_status, last_health_check_at, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name, type = EXCLUDED.type, base_url = EXCLUDED.base_url,
			api_key_env = EXCLUDED.api_key_env,
			api_key = CASE WHEN EXCLUDED.api_key = '' THEN ai.providers.api_key ELSE EXCLUDED.api_key END,
			health_status = EXCLUDED.health_status,
			last_health_check_at = EXCLUDED.last_health_check_at`,
		p.ID, p.Name, p.Type, p.BaseURL, p.APIKeyEnv, p.APIKey, p.HealthStatus, p.LastHealthCheckAt, p.CreatedAt)
	return err
}

func scanProvider(sc interface {
	Scan(dest ...interface{}) error
}) (*domain.Provider, error) {
	var p domain.Provider
	err := sc.Scan(&p.ID, &p.Name, &p.Type, &p.BaseURL, &p.APIKeyEnv, &p.APIKey,
		&p.HealthStatus, &p.LastHealthCheckAt, &p.CreatedAt)
	return &p, err
}

func (r *PostgresModelRegistry) FindProviderByID(ctx context.Context, id domain.ProviderID) (*domain.Provider, error) {
	p, err := scanProvider(r.q.QueryRow(ctx, `SELECT id, name, type, base_url, api_key_env, api_key, health_status, last_health_check_at, created_at FROM ai.providers WHERE id = $1`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrModelNotFound
	}
	return p, err
}

func (r *PostgresModelRegistry) ListProviders(ctx context.Context) ([]*domain.Provider, error) {
	rows, err := r.q.Query(ctx, `SELECT id, name, type, base_url, api_key_env, api_key, health_status, last_health_check_at, created_at FROM ai.providers ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []*domain.Provider{}
	for rows.Next() {
		p, err := scanProvider(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *PostgresModelRegistry) DeleteProvider(ctx context.Context, id domain.ProviderID) error {
	if _, err := r.q.Exec(ctx, `DELETE FROM ai.models WHERE provider_id = $1`, id); err != nil {
		return err
	}
	_, err := r.q.Exec(ctx, `DELETE FROM ai.providers WHERE id = $1`, id)
	return err
}

// --- Models ---

func (r *PostgresModelRegistry) SaveModel(ctx context.Context, m *domain.Model) error {
	caps := make([]string, len(m.Capabilities))
	for i, c := range m.Capabilities {
		caps[i] = string(c)
	}
	if m.Status == "" {
		m.Status = domain.ModelActive
	}
	_, err := r.q.Exec(ctx, `
		INSERT INTO ai.models (id, provider_id, name, identifier, capabilities, modalities, supported_parameters, limits, pricing, version, status, metadata, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		ON CONFLICT (provider_id, identifier) DO UPDATE SET
			name = EXCLUDED.name, capabilities = EXCLUDED.capabilities,
			modalities = EXCLUDED.modalities, supported_parameters = EXCLUDED.supported_parameters,
			limits = EXCLUDED.limits, pricing = EXCLUDED.pricing,
			version = EXCLUDED.version, status = EXCLUDED.status, metadata = EXCLUDED.metadata`,
		m.ID, m.ProviderID, m.Name, m.Identifier, mustJSON(caps), mustJSON(m.Modalities),
		mustJSON(m.SupportedParameters), mustJSON(m.Limits), mustJSON(m.Pricing),
		m.Version, string(m.Status), mustJSON(m.Metadata), m.CreatedAt)
	return err
}

func scanModel(sc interface {
	Scan(dest ...interface{}) error
}) (*domain.Model, error) {
	var m domain.Model
	var caps, modalities, params, limits, pricing, meta []byte
	err := sc.Scan(&m.ID, &m.ProviderID, &m.Name, &m.Identifier, &caps, &modalities,
		&params, &limits, &pricing, &m.Version, &m.Status, &meta, &m.CreatedAt)
	var names []string
	_ = json.Unmarshal(caps, &names)
	for _, n := range names {
		m.Capabilities = append(m.Capabilities, domain.AICapability(n))
	}
	_ = json.Unmarshal(modalities, &m.Modalities)
	_ = json.Unmarshal(params, &m.SupportedParameters)
	_ = json.Unmarshal(limits, &m.Limits)
	_ = json.Unmarshal(pricing, &m.Pricing)
	_ = json.Unmarshal(meta, &m.Metadata)
	return &m, err
}

const modelCols = `id, provider_id, name, identifier, capabilities, modalities, supported_parameters, limits, pricing, version, status, metadata, created_at`

func (r *PostgresModelRegistry) FindModelByID(ctx context.Context, id domain.ModelID) (*domain.Model, error) {
	m, err := scanModel(r.q.QueryRow(ctx, `SELECT id, provider_id, name, identifier, capabilities, modalities, supported_parameters, limits, pricing, version, status, metadata, created_at FROM ai.models WHERE id = $1`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrModelNotFound
	}
	return m, err
}

func (r *PostgresModelRegistry) ListModels(ctx context.Context, providerID domain.ProviderID) ([]*domain.Model, error) {
	rows, err := r.q.Query(ctx, `SELECT id, provider_id, name, identifier, capabilities, modalities, supported_parameters, limits, pricing, version, status, metadata, created_at FROM ai.models WHERE $1::text = '' OR provider_id = $1 ORDER BY name`, string(providerID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []*domain.Model{}
	for rows.Next() {
		m, err := scanModel(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (r *PostgresModelRegistry) ListModelsByCapability(ctx context.Context, capability domain.AICapability) ([]*domain.Model, error) {
	rows, err := r.q.Query(ctx, `SELECT id, provider_id, name, identifier, capabilities, modalities, supported_parameters, limits, pricing, version, status, metadata, created_at FROM ai.models WHERE capabilities @> to_jsonb($1::text) ORDER BY name`, string(capability))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []*domain.Model{}
	for rows.Next() {
		m, err := scanModel(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (r *PostgresModelRegistry) DeleteModel(ctx context.Context, id domain.ModelID) error {
	_, err := r.q.Exec(ctx, `DELETE FROM ai.models WHERE id = $1`, id)
	return err
}

// --- Policies ---

func (r *PostgresModelRegistry) UpsertPolicy(ctx context.Context, p *domain.PolicyRow) error {
	_, err := r.q.Exec(ctx, `
		INSERT INTO ai.model_policies (id, scope, scope_id, capability, provider_id, model_id, fallback_models, allowed_models, preferred_models, routing_strategy, quality_requirement, max_cost, max_latency_ms, region, concurrency_limit, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
		ON CONFLICT (scope, scope_id, capability) DO UPDATE SET
			provider_id = EXCLUDED.provider_id, model_id = EXCLUDED.model_id,
			fallback_models = EXCLUDED.fallback_models, allowed_models = EXCLUDED.allowed_models,
			preferred_models = EXCLUDED.preferred_models, routing_strategy = EXCLUDED.routing_strategy,
			quality_requirement = EXCLUDED.quality_requirement, max_cost = EXCLUDED.max_cost,
			max_latency_ms = EXCLUDED.max_latency_ms, region = EXCLUDED.region,
			concurrency_limit = EXCLUDED.concurrency_limit`,
		p.ID, p.Scope, p.ScopeID, string(p.Capability), p.ProviderID, p.ModelID,
		mustJSON(p.FallbackModels), mustJSON(p.AllowedModels), mustJSON(p.PreferredModels),
		p.Strategy(), p.QualityRequirement, p.MaxCost, p.MaxLatencyMs, p.Region,
		p.ConcurrencyLimit, p.CreatedAt)
	return err
}

func scanPolicy(sc interface {
	Scan(dest ...interface{}) error
}) (*domain.PolicyRow, error) {
	var p domain.PolicyRow
	var fallback, allowed, preferred []byte
	err := sc.Scan(&p.ID, &p.Scope, &p.ScopeID, &p.Capability, &p.ProviderID, &p.ModelID,
		&fallback, &allowed, &preferred, &p.RoutingStrategy, &p.QualityRequirement,
		&p.MaxCost, &p.MaxLatencyMs, &p.Region, &p.ConcurrencyLimit, &p.CreatedAt)
	_ = json.Unmarshal(fallback, &p.FallbackModels)
	_ = json.Unmarshal(allowed, &p.AllowedModels)
	_ = json.Unmarshal(preferred, &p.PreferredModels)
	return &p, err
}

const policyCols = `id, scope, scope_id, capability, provider_id, model_id, fallback_models, allowed_models, preferred_models, routing_strategy, quality_requirement, max_cost, max_latency_ms, region, concurrency_limit, created_at`

// ResolvePolicy returns the policy row for the exact scope id passed. The
// scope string encodes level:id, e.g. "project:abc" or "system".
func (r *PostgresModelRegistry) ResolvePolicy(ctx context.Context, capability, scopeID string) (*domain.PolicyRow, error) {
	scope, id := splitScope(scopeID)
	p, err := scanPolicy(r.q.QueryRow(ctx, `SELECT id, scope, scope_id, capability, provider_id, model_id, fallback_models, allowed_models, preferred_models, routing_strategy, quality_requirement, max_cost, max_latency_ms, region, concurrency_limit, created_at FROM ai.model_policies WHERE scope = $1 AND scope_id = $2 AND capability = $3`, scope, id, capability))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil // no policy at this scope; caller walks the chain
	}
	return p, err
}

func splitScope(scopeID string) (scope, id string) {
	for i, c := range scopeID {
		if c == ':' {
			return scopeID[:i], scopeID[i+1:]
		}
	}
	if scopeID == "" {
		return domain.ScopeSystem, ""
	}
	return scopeID, ""
}

func (r *PostgresModelRegistry) ListPolicies(ctx context.Context, scope, scopeID string) ([]*domain.PolicyRow, error) {
	rows, err := r.q.Query(ctx, `SELECT id, scope, scope_id, capability, provider_id, model_id, fallback_models, allowed_models, preferred_models, routing_strategy, quality_requirement, max_cost, max_latency_ms, region, concurrency_limit, created_at FROM ai.model_policies WHERE scope = $1 AND scope_id = $2 ORDER BY capability`, scope, scopeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []*domain.PolicyRow{}
	for rows.Next() {
		p, err := scanPolicy(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// --- Generation jobs ---

func scanGenJob(sc interface {
	Scan(dest ...interface{}) error
}) (*domain.GenerationJob, error) {
	var j domain.GenerationJob
	var key *string
	err := sc.Scan(&j.ID, &j.ProjectID, &j.Capability, &j.ProviderID, &j.ModelID, &j.ModelVersion,
		&j.Input, &j.Output, &j.Status, &j.Attempt, &j.Cost, &j.ProviderJobID,
		&key, &j.Error, &j.StartedAt, &j.CompletedAt, &j.CreatedAt)
	if key != nil {
		j.IdempotencyKey = *key
	}
	return &j, err
}

const genJobCols = `id, project_id, capability, provider_id, model_id, model_version, input, output, status, attempt, cost, provider_job_id, idempotency_key, error, started_at, completed_at, created_at`

func (r *PostgresModelRegistry) SaveGenerationJob(ctx context.Context, j *domain.GenerationJob) error {
	var key *string
	if j.IdempotencyKey != "" {
		key = &j.IdempotencyKey
	}
	status := string(j.Status)
	if status == "" {
		status = string(domain.StatusPending)
	}
	input := j.Input
	if len(input) == 0 {
		input = json.RawMessage(`{}`)
	}
	output := j.Output
	if len(output) == 0 {
		output = json.RawMessage(`{}`)
	}
	_, err := r.q.Exec(ctx, `
		INSERT INTO ai.generation_jobs (id, project_id, capability, provider_id, model_id, model_version, input, output, status, attempt, cost, provider_job_id, idempotency_key, error, started_at, completed_at, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
		ON CONFLICT (id) DO UPDATE SET
			output = EXCLUDED.output, status = EXCLUDED.status,
			attempt = EXCLUDED.attempt, cost = EXCLUDED.cost,
			provider_job_id = EXCLUDED.provider_job_id, error = EXCLUDED.error,
			started_at = EXCLUDED.started_at, completed_at = EXCLUDED.completed_at`,
		j.ID, j.ProjectID, string(j.Capability), j.ProviderID, j.ModelID, j.ModelVersion,
		input, output, status, j.Attempt, j.Cost, j.ProviderJobID, key,
		j.Error, j.StartedAt, j.CompletedAt, j.CreatedAt)
	return err
}

func (r *PostgresModelRegistry) FindGenerationJobByID(ctx context.Context, id string) (*domain.GenerationJob, error) {
	j, err := scanGenJob(r.q.QueryRow(ctx, `SELECT id, project_id, capability, provider_id, model_id, model_version, input, output, status, attempt, cost, provider_job_id, idempotency_key, error, started_at, completed_at, created_at FROM ai.generation_jobs WHERE id = $1`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrModelNotFound
	}
	return j, err
}

func (r *PostgresModelRegistry) FindGenerationJobByIdempotencyKey(ctx context.Context, key string) (*domain.GenerationJob, error) {
	j, err := scanGenJob(r.q.QueryRow(ctx, `SELECT id, project_id, capability, provider_id, model_id, model_version, input, output, status, attempt, cost, provider_job_id, idempotency_key, error, started_at, completed_at, created_at FROM ai.generation_jobs WHERE idempotency_key = $1`, key))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrModelNotFound
	}
	return j, err
}

func (r *PostgresModelRegistry) ListGenerationJobs(ctx context.Context, projectID string) ([]*domain.GenerationJob, error) {
	rows, err := r.q.Query(ctx, `SELECT id, project_id, capability, provider_id, model_id, model_version, input, output, status, attempt, cost, provider_job_id, idempotency_key, error, started_at, completed_at, created_at FROM ai.generation_jobs WHERE project_id = $1 ORDER BY created_at DESC`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []*domain.GenerationJob{}
	for rows.Next() {
		j, err := scanGenJob(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, j)
	}
	return out, rows.Err()
}
