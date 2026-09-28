package persistence

import (
	"context"
	"encoding/json"

	"dramastudio/internal/intelligence/domain"
	"dramastudio/internal/platform/database/postgres"
)

// PostgresRepository persists execution records and policy layers. It runs
// on the postgres.Querier port — SQLite serves the same statements.
type PostgresRepository struct {
	q postgres.Querier
}

func NewPostgresRepository(q postgres.Querier) *PostgresRepository {
	return &PostgresRepository{q: q}
}

func (r *PostgresRepository) SaveExecution(ctx context.Context, rec *domain.ExecutionRecord) error {
	viol, _ := json.Marshal(rec.Violations)
	find, _ := json.Marshal(rec.Findings)
	_, err := r.q.Exec(ctx, `
		INSERT INTO intelligence.executions
		  (id, project_id, task_id, action, agent_id, agent_version,
		   skill_ids, skill_versions, policy_ids, rule_ids, evaluator_ids,
		   provider_id, model_id, model_version, prompt_hash, context_hash,
		   generation_job_id, verdict, violations, findings, cost, error, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23)
		ON CONFLICT (id) DO UPDATE SET verdict = EXCLUDED.verdict,
			violations = EXCLUDED.violations, findings = EXCLUDED.findings,
			error = EXCLUDED.error`,
		rec.ID, rec.ProjectID, rec.TaskID, rec.Action, rec.AgentID, rec.AgentVersion,
		mustJSON(rec.SkillIDs), mustJSON(rec.SkillVersions), mustJSON(rec.PolicyIDs),
		mustJSON(rec.RuleIDs), mustJSON(rec.EvaluatorIDs),
		rec.ProviderID, rec.ModelID, rec.ModelVersion, rec.PromptHash, rec.ContextHash,
		rec.GenerationJobID, string(rec.Verdict), viol, find, rec.Cost, rec.Error, rec.CreatedAt)
	return err
}

func (r *PostgresRepository) FindExecution(ctx context.Context, id string) (*domain.ExecutionRecord, error) {
	row := r.q.QueryRow(ctx, `
		SELECT id, project_id, task_id, action, agent_id, agent_version,
		       skill_ids, skill_versions, policy_ids, rule_ids, evaluator_ids,
		       provider_id, model_id, model_version, prompt_hash, context_hash,
		       generation_job_id, verdict, violations, findings, cost, error, created_at
		FROM intelligence.executions WHERE id = $1`, id)
	rec, err := scanExecution(row)
	if err != nil {
		if postgres.IsNoRows(err) {
			return nil, domain.ErrExecutionNotFound
		}
		return nil, err
	}
	return rec, nil
}

func (r *PostgresRepository) ListExecutions(ctx context.Context, projectID string) ([]*domain.ExecutionRecord, error) {
	rows, err := r.q.Query(ctx, `
		SELECT id, project_id, task_id, action, agent_id, agent_version,
		       skill_ids, skill_versions, policy_ids, rule_ids, evaluator_ids,
		       provider_id, model_id, model_version, prompt_hash, context_hash,
		       generation_job_id, verdict, violations, findings, cost, error, created_at
		FROM intelligence.executions
		WHERE $1::text = '' OR project_id = $1
		ORDER BY created_at DESC`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.ExecutionRecord
	for rows.Next() {
		rec, err := scanExecution(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, rec)
	}
	return out, rows.Err()
}

func (r *PostgresRepository) SavePolicyLayer(ctx context.Context, p *domain.Policy) error {
	_, err := r.q.Exec(ctx, `
		INSERT INTO intelligence.policy_layers
		  (policy_id, layer, scope_id, version, applies_to, requirements, protected, constraints)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		ON CONFLICT (policy_id, layer, scope_id) DO UPDATE SET
			version = EXCLUDED.version, applies_to = EXCLUDED.applies_to,
			requirements = EXCLUDED.requirements, protected = EXCLUDED.protected,
			constraints = EXCLUDED.constraints`,
		p.ID, string(p.Layer), p.ScopeID, p.Version,
		mustJSON(p.AppliesTo), mustJSON(p.Requirements), mustJSON(p.Protected), mustJSON(p.Constraints))
	return err
}

func (r *PostgresRepository) ListPolicyLayers(ctx context.Context, policyID string) ([]*domain.Policy, error) {
	rows, err := r.q.Query(ctx, `
		SELECT policy_id, layer, scope_id, version, applies_to, requirements, protected, constraints
		FROM intelligence.policy_layers WHERE policy_id = $1`, policyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.Policy
	for rows.Next() {
		p := &domain.Policy{}
		var layer string
		var appliesTo, reqs, prot, cons []byte
		if err := rows.Scan(&p.ID, &layer, &p.ScopeID, &p.Version, &appliesTo, &reqs, &prot, &cons); err != nil {
			return nil, err
		}
		p.Layer = domain.PolicyLayer(layer)
		_ = json.Unmarshal(appliesTo, &p.AppliesTo)
		_ = json.Unmarshal(reqs, &p.Requirements)
		_ = json.Unmarshal(prot, &p.Protected)
		_ = json.Unmarshal(cons, &p.Constraints)
		out = append(out, p)
	}
	return out, rows.Err()
}

type scanner interface {
	Scan(dest ...interface{}) error
}

func scanExecution(row scanner) (*domain.ExecutionRecord, error) {
	rec := &domain.ExecutionRecord{}
	var skills, svers, pols, rules, evals, viol, find []byte
	var verdict string
	err := row.Scan(&rec.ID, &rec.ProjectID, &rec.TaskID, &rec.Action, &rec.AgentID, &rec.AgentVersion,
		&skills, &svers, &pols, &rules, &evals,
		&rec.ProviderID, &rec.ModelID, &rec.ModelVersion, &rec.PromptHash, &rec.ContextHash,
		&rec.GenerationJobID, &verdict, &viol, &find, &rec.Cost, &rec.Error, &rec.CreatedAt)
	if err != nil {
		return nil, err
	}
	rec.Verdict = domain.ExecutionVerdict(verdict)
	_ = json.Unmarshal(skills, &rec.SkillIDs)
	_ = json.Unmarshal(svers, &rec.SkillVersions)
	_ = json.Unmarshal(pols, &rec.PolicyIDs)
	_ = json.Unmarshal(rules, &rec.RuleIDs)
	_ = json.Unmarshal(evals, &rec.EvaluatorIDs)
	_ = json.Unmarshal(viol, &rec.Violations)
	_ = json.Unmarshal(find, &rec.Findings)
	return rec, nil
}

func mustJSON(v interface{}) []byte {
	b, _ := json.Marshal(v)
	if len(b) == 0 || string(b) == "null" {
		return []byte("[]")
	}
	return b
}
