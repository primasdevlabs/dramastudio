package persistence

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"

	"dramastudio/internal/agents/domain"
	"dramastudio/internal/platform/database/postgres"
)

type PostgresAgentRepository struct {
	q postgres.Querier
}

func NewPostgresAgentRepository(q postgres.Querier) *PostgresAgentRepository {
	return &PostgresAgentRepository{q: q}
}

func mustJSON(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}

// --- Definitions ---

func (r *PostgresAgentRepository) SaveDefinition(ctx context.Context, d *domain.Definition) error {
	_, err := r.q.Exec(ctx, `
		INSERT INTO agents.definitions (id, role, name, instructions, skills, tools, permissions, model_policy, budget_limit, active)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name, instructions = EXCLUDED.instructions,
			skills = EXCLUDED.skills, tools = EXCLUDED.tools,
			permissions = EXCLUDED.permissions, model_policy = EXCLUDED.model_policy,
			budget_limit = EXCLUDED.budget_limit, active = EXCLUDED.active`,
		d.ID, string(d.Role), d.Name, d.Instructions, mustJSON(d.Skills),
		mustJSON(d.Tools), mustJSON(d.Permissions), mustJSON(d.ModelPolicy),
		d.BudgetLimit, d.Active)
	return err
}

func scanDefinition(sc interface {
	Scan(dest ...interface{}) error
}) (*domain.Definition, error) {
	var d domain.Definition
	var skills, tools, perms, policy []byte
	err := sc.Scan(&d.ID, &d.Role, &d.Name, &d.Instructions, &skills, &tools,
		&perms, &policy, &d.BudgetLimit, &d.Active)
	_ = json.Unmarshal(skills, &d.Skills)
	_ = json.Unmarshal(tools, &d.Tools)
	_ = json.Unmarshal(perms, &d.Permissions)
	_ = json.Unmarshal(policy, &d.ModelPolicy)
	return &d, err
}

const defColumns = `id, role, name, instructions, skills, tools, permissions, model_policy, budget_limit, active`

func (r *PostgresAgentRepository) FindDefinitionByID(ctx context.Context, id string) (*domain.Definition, error) {
	d, err := scanDefinition(r.q.QueryRow(ctx, `SELECT id, role, name, instructions, skills, tools, permissions, model_policy, budget_limit, active FROM agents.definitions WHERE id = $1`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrDefinitionNotFound
	}
	return d, err
}

func (r *PostgresAgentRepository) ListDefinitions(ctx context.Context) ([]*domain.Definition, error) {
	rows, err := r.q.Query(ctx, `SELECT id, role, name, instructions, skills, tools, permissions, model_policy, budget_limit, active FROM agents.definitions ORDER BY role`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []*domain.Definition{}
	for rows.Next() {
		d, err := scanDefinition(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

// --- Tasks ---

const taskColumns = `id, project_id, agent_id, objective, input, expected_output, constraints,
	budget, max_iterations, iteration, timeout_sec, approval_policy, success_criteria,
	status, result, error, created_at, completed_at`

func scanTask(sc interface {
	Scan(dest ...interface{}) error
}) (*domain.Task, error) {
	var t domain.Task
	var input, constraints, result []byte
	var policy, status string
	err := sc.Scan(&t.ID, &t.ProjectID, &t.AgentID, &t.Objective, &input,
		&t.ExpectedOutput, &constraints, &t.Budget, &t.MaxIterations, &t.Iteration,
		&t.TimeoutSec, &policy, &t.SuccessCriteria, &status, &result, &t.Error,
		&t.CreatedAt, &t.CompletedAt)
	t.ApprovalPolicy = domain.ApprovalPolicy(policy)
	t.Status = domain.TaskStatus(status)
	_ = json.Unmarshal(input, &t.Input)
	_ = json.Unmarshal(constraints, &t.Constraints)
	_ = json.Unmarshal(result, &t.Result)
	return &t, err
}

func (r *PostgresAgentRepository) SaveTask(ctx context.Context, t *domain.Task) error {
	status := string(t.Status)
	if status == "" {
		status = string(domain.TaskPending)
	}
	policy := string(t.ApprovalPolicy)
	if policy == "" {
		policy = string(domain.PolicyMonitored)
	}
	maxIter := t.MaxIterations
	if maxIter == 0 {
		maxIter = 3
	}
	timeout := t.TimeoutSec
	if timeout == 0 {
		timeout = 3600
	}
	_, err := r.q.Exec(ctx, `
		INSERT INTO agents.tasks (id, project_id, agent_id, objective, input, expected_output, constraints,
			budget, max_iterations, iteration, timeout_sec, approval_policy, success_criteria,
			status, result, error, created_at, completed_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18)
		ON CONFLICT (id) DO UPDATE SET
			iteration = EXCLUDED.iteration, status = EXCLUDED.status,
			result = EXCLUDED.result, error = EXCLUDED.error,
			completed_at = EXCLUDED.completed_at`,
		t.ID, t.ProjectID, t.AgentID, t.Objective, mustJSON(t.Input),
		t.ExpectedOutput, mustJSON(t.Constraints), t.Budget, maxIter, t.Iteration,
		timeout, policy, t.SuccessCriteria, status, nullableJSON(t.Result), t.Error,
		t.CreatedAt, t.CompletedAt)
	return err
}

func nullableJSON(m map[string]interface{}) interface{} {
	if m == nil {
		return nil
	}
	return mustJSON(m)
}

func (r *PostgresAgentRepository) FindTaskByID(ctx context.Context, id string) (*domain.Task, error) {
	t, err := scanTask(r.q.QueryRow(ctx, `SELECT id, project_id, agent_id, objective, input, expected_output, constraints,
	budget, max_iterations, iteration, timeout_sec, approval_policy, success_criteria,
	status, result, error, created_at, completed_at FROM agents.tasks WHERE id = $1`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrTaskNotFound
	}
	return t, err
}

func (r *PostgresAgentRepository) ListTasks(ctx context.Context, projectID string, status domain.TaskStatus) ([]*domain.Task, error) {
	sql := `SELECT id, project_id, agent_id, objective, input, expected_output, constraints,
	budget, max_iterations, iteration, timeout_sec, approval_policy, success_criteria,
	status, result, error, created_at, completed_at FROM agents.tasks WHERE project_id = $1`
	args := []interface{}{projectID}
	if status != "" {
		sql += ` AND status = $2`
		args = append(args, string(status))
	}
	sql += ` ORDER BY created_at DESC`
	rows, err := r.q.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []*domain.Task{}
	for rows.Next() {
		t, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

// --- Decisions ---

func (r *PostgresAgentRepository) SaveDecision(ctx context.Context, d *domain.Decision) error {
	_, err := r.q.Exec(ctx, `
		INSERT INTO agents.decisions (id, project_id, episode_id, decision, reason, decision_maker, mode, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		d.ID, d.ProjectID, d.EpisodeID, d.Decision, d.Reason,
		string(d.DecisionMaker), d.Mode, d.Timestamp)
	return err
}

func (r *PostgresAgentRepository) ListDecisionsByProject(ctx context.Context, projectID string) ([]*domain.Decision, error) {
	rows, err := r.q.Query(ctx, `
		SELECT id, project_id, episode_id, decision, reason, decision_maker, mode, created_at
		FROM agents.decisions WHERE project_id = $1 ORDER BY created_at DESC`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []*domain.Decision{}
	for rows.Next() {
		var d domain.Decision
		if err := rows.Scan(&d.ID, &d.ProjectID, &d.EpisodeID, &d.Decision,
			&d.Reason, &d.DecisionMaker, &d.Mode, &d.Timestamp); err != nil {
			return nil, err
		}
		out = append(out, &d)
	}
	return out, rows.Err()
}
