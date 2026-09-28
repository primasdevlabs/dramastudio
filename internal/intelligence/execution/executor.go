// Package execution is the runtime the production pipeline calls:
//
//	Task → Resolve Agent → Build ExecutionContext → Resolve Model →
//	Execute → Validate Rules → Evaluate Quality → Commit Result
//
// The model never overrides the layers above it: rules are machine-enforced,
// evaluators grade independently, and every run writes a versioned
// ExecutionRecord for reproducibility.
package execution

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"dramastudio/internal/intelligence/domain"
	"dramastudio/internal/intelligence/evaluators"
	"dramastudio/internal/intelligence/resolver"
	"dramastudio/internal/intelligence/rules"
)

// GenResult is the provider-agnostic output of one generation dispatch.
type GenResult struct {
	JobID        string          `json:"job_id"`
	ProviderID   string          `json:"provider_id"`
	ModelID      string          `json:"model_id"`
	ModelVersion string          `json:"model_version"`
	Output       json.RawMessage `json:"output"`
	Cost         float64         `json:"cost"`
	Status       string          `json:"status"` // succeeded|submitted (async)
}

// Generator is the seam into the AI subsystem — wired to
// ai.ExecuteGenerationHandler so model routing, fallbacks, and provenance
// stay unified.
type Generator interface {
	Generate(ctx context.Context, task domain.TaskContext, prompt string) (*GenResult, error)
}

// Executor resolves and runs one task against the intelligence layer.
type Executor struct {
	reg   *resolver.Registry
	rules *rules.Engine
	evals *evaluators.Runner
	gen   Generator
	repo  domain.ExecutionRepository
}

func NewExecutor(reg *resolver.Registry, e *evaluators.Runner, gen Generator, repo domain.ExecutionRepository) *Executor {
	return &Executor{reg: reg, rules: rules.NewEngine(), evals: e, gen: gen, repo: repo}
}

// Run executes the task end to end and records an ExecutionRecord.
func (x *Executor) Run(ctx context.Context, agentID string, task domain.TaskContext) (*domain.ExecutionRecord, *GenResult, error) {
	ec, err := x.reg.ResolveExecution(ctx, agentID, task)
	if err != nil {
		return nil, nil, err
	}
	rec := x.newRecord(ec)

	// 1. Pre-generation rules — BLOCK stops before any provider call.
	pre := x.rules.Check(ec.Rules, task.Action, "pre", task.Facts)
	pre = append(pre, x.rules.Check(ec.Rules, task.Action, "", task.Facts)...)
	rec.Violations = append(rec.Violations, pre...)
	if v := firstOf(pre, domain.EnforceBlock); v != nil {
		rec.Verdict = domain.VerdictBlocked
		rec.Error = v.Message
		return rec, nil, x.save(ctx, rec)
	}

	// 2. Generate — prompt carries skill method + policy constraints, not
	//    hardcoded workflow text.
	prompt := buildPrompt(task.Prompt, ec.InstructionLines)
	res, err := x.gen.Generate(ctx, task, prompt)
	if err != nil {
		rec.Verdict = domain.VerdictFailed
		rec.Error = err.Error()
		return rec, nil, x.save(ctx, rec)
	}
	rec.GenerationJobID = res.JobID
	rec.ProviderID, rec.ModelID, rec.ModelVersion = res.ProviderID, res.ModelID, res.ModelVersion
	rec.Cost = res.Cost
	rec.PromptHash = hash(prompt)

	// 3. Post-generation rules over task facts + output.
	outFacts := mergeFacts(task.Facts, res.Output)
	post := x.rules.Check(ec.Rules, task.Action, "post", outFacts)
	rec.Violations = append(rec.Violations, post...)

	// 4. Independent evaluation.
	for _, ev := range ec.Evaluators {
		findings, ferr := x.evals.Run(ctx, ev, task.Action, outFacts)
		if ferr != nil {
			rec.Findings = append(rec.Findings, domain.EvalFinding{
				EvaluatorID: ev.ID, Check: "runner", Severity: domain.SeverityWarning,
				Detail: ferr.Error(),
			})
			continue
		}
		rec.Findings = append(rec.Findings, findings...)
	}

	rec.Verdict = verdict(rec)
	return rec, res, x.save(ctx, rec)
}

// Preview resolves the context without executing — the "what would this task
// see" inspection endpoint.
func (x *Executor) Preview(ctx context.Context, agentID string, task domain.TaskContext) (*domain.AgentExecutionContext, error) {
	return x.reg.ResolveExecution(ctx, agentID, task)
}

func (x *Executor) newRecord(ec *domain.AgentExecutionContext) *domain.ExecutionRecord {
	rec := &domain.ExecutionRecord{
		ID: "exec_" + uuid.NewString(), ProjectID: ec.Task.ProjectID,
		TaskID: ec.Task.ID, Action: ec.Task.Action,
		AgentID: ec.Agent.ID, AgentVersion: ec.Agent.Version,
		CreatedAt: time.Now().UTC(),
	}
	for _, s := range ec.Skills {
		rec.SkillIDs = append(rec.SkillIDs, s.ID)
		rec.SkillVersions = append(rec.SkillVersions, s.Version)
	}
	for _, p := range ec.Policies {
		rec.PolicyIDs = append(rec.PolicyIDs, p.PolicyID)
	}
	for _, rl := range ec.Rules {
		rec.RuleIDs = append(rec.RuleIDs, rl.ID)
	}
	for _, ev := range ec.Evaluators {
		rec.EvaluatorIDs = append(rec.EvaluatorIDs, ev.ID)
	}
	rec.ContextHash = hash(ec.Task.Action + ":" + string(ec.Task.Input) + ":" + fmt.Sprint(ec.Task.Facts))
	return rec
}

func (x *Executor) save(ctx context.Context, rec *domain.ExecutionRecord) error {
	if x.repo == nil {
		return nil
	}
	return x.repo.SaveExecution(ctx, rec)
}

// verdict picks the strictest outcome across violations and findings.
func verdict(rec *domain.ExecutionRecord) domain.ExecutionVerdict {
	v := domain.VerdictApproved
	for _, viol := range rec.Violations {
		switch viol.Enforcement {
		case domain.EnforceBlock:
			return domain.VerdictBlocked
		case domain.EnforceRequireReview:
			v = domain.VerdictNeedsReview
		case domain.EnforceAutoFix:
			if v == domain.VerdictApproved {
				v = domain.VerdictAutoFixMarked
			}
		}
	}
	for _, f := range rec.Findings {
		if f.Severity == domain.SeverityBlocking {
			return domain.VerdictBlocked
		}
		if f.Severity == domain.SeverityError && v == domain.VerdictApproved {
			v = domain.VerdictNeedsReview
		}
	}
	return v
}

func firstOf(vs []domain.RuleViolation, e domain.Enforcement) *domain.RuleViolation {
	for i := range vs {
		if vs[i].Enforcement == e {
			return &vs[i]
		}
	}
	return nil
}

// buildPrompt assembles task instructions + resolved skill methods + policy
// constraints — the dynamic replacement for workflow-embedded prompt text.
func buildPrompt(base string, lines []string) string {
	if len(lines) == 0 {
		return base
	}
	out := base + "\n\nProduction constraints and methodology:\n"
	for _, l := range lines {
		out += "- " + l + "\n"
	}
	return out
}

// mergeFacts overlays the generation output under "output" for post rules.
func mergeFacts(facts map[string]interface{}, output json.RawMessage) map[string]interface{} {
	merged := map[string]interface{}{}
	for k, v := range facts {
		merged[k] = v
	}
	var out interface{}
	if err := json.Unmarshal(output, &out); err == nil {
		merged["output"] = out
	}
	return merged
}

func hash(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:16])
}
