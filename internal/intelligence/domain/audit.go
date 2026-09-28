package domain

import "time"

// RuleViolation is one failed predicate from the rules engine.
type RuleViolation struct {
	RuleID      string      `json:"rule_id"`
	RuleVersion int         `json:"rule_version"`
	Predicate   string      `json:"predicate"`
	Severity    Severity    `json:"severity"`
	Enforcement Enforcement `json:"enforcement"`
	Message     string      `json:"message"`
	Phase       string      `json:"phase"` // pre|post
}

// EvalFinding is one evaluator observation.
type EvalFinding struct {
	EvaluatorID string   `json:"evaluator_id"`
	Check       string   `json:"check"`
	Severity    Severity `json:"severity"`
	Detail      string   `json:"detail"`
}

// ExecutionVerdict is the gate result for a task execution.
type ExecutionVerdict string

const (
	VerdictApproved      ExecutionVerdict = "approved"       // clean
	VerdictNeedsReview   ExecutionVerdict = "needs_review"   // REQUIRE_REVIEW violations
	VerdictBlocked       ExecutionVerdict = "blocked"        // BLOCK violations
	VerdictFailed        ExecutionVerdict = "failed"         // generation failed
	VerdictAutoFixMarked ExecutionVerdict = "autofix_marked" // AUTO_FIX violations
)

// ExecutionRecord is the reproducibility audit for one task execution —
// every versioned input that shaped the output (§9 of the intelligence spec:
// episode 40 must be reproducible even though skills moved to v3).
type ExecutionRecord struct {
	ID        string `json:"id"`
	ProjectID string `json:"project_id"`
	TaskID    string `json:"task_id"`
	Action    string `json:"action"`

	AgentID       string   `json:"agent_id"`
	AgentVersion  int      `json:"agent_version"`
	SkillIDs      []string `json:"skill_ids"`
	SkillVersions []int    `json:"skill_versions"`
	PolicyIDs     []string `json:"policy_ids"`
	RuleIDs       []string `json:"rule_ids"`
	EvaluatorIDs  []string `json:"evaluator_ids"`

	ProviderID   string `json:"provider_id"`
	ModelID      string `json:"model_id"`
	ModelVersion string `json:"model_version"`

	PromptHash  string `json:"prompt_hash"`
	ContextHash string `json:"context_hash"`

	GenerationJobID string           `json:"generation_job_id,omitempty"`
	Verdict         ExecutionVerdict `json:"verdict"`
	Violations      []RuleViolation  `json:"violations,omitempty"`
	Findings        []EvalFinding    `json:"findings,omitempty"`
	Cost            float64          `json:"cost"`
	Error           string           `json:"error,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}
