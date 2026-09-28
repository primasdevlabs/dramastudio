package domain

import "encoding/json"

// TaskContext is the unit of work the executor resolves against — the
// pipeline's "Create Task → Resolve Agent" bridge.
type TaskContext struct {
	ID         string                 `json:"id"`
	ProjectID  string                 `json:"project_id"`
	Scope      string                 `json:"scope"`      // system|studio|project|series:x|season:x|episode:x|task:x
	Action     string                 `json:"action"`     // e.g. generate_shot, write_episode_script
	Capability string                 `json:"capability"` // AI capability this task routes to
	Prompt     string                 `json:"prompt"`
	Input      json.RawMessage        `json:"input,omitempty"`
	Facts      map[string]interface{} `json:"facts,omitempty"` // rule-evaluation fact bag
	Metadata   map[string]string      `json:"metadata,omitempty"`
}

// AgentExecutionContext is the resolved bundle handed to the executor —
// exactly what the agent may see and the constraints that apply. Built by
// the resolver, consumed by the executor, recorded by the audit trail.
type AgentExecutionContext struct {
	Agent      *Agent             `json:"agent"`
	Skills     []*Skill           `json:"skills"`
	Policies   []*EffectivePolicy `json:"policies"`
	Rules      []*Rule            `json:"rules"`
	Evaluators []*Evaluator       `json:"evaluators"`

	Task        TaskContext   `json:"task"`
	Permissions PermissionSet `json:"permissions"`

	// InstructionLines is the assembled method + constraint text injected
	// into the model call (skill methods then policy constraints).
	InstructionLines []string `json:"instruction_lines"`
}

// EffectivePolicy is a policy after layered resolution — the merged
// requirements plus where each winning value came from.
type EffectivePolicy struct {
	PolicyID     string                 `json:"policy_id"`
	Requirements map[string]interface{} `json:"requirements"`
	Constraints  []string               `json:"constraints"`
	// Trace records which layer produced each requirement key (audit).
	Trace map[string]string `json:"trace,omitempty"`
	// Dropped lists lower-layer overrides rejected by protected keys.
	Dropped []string `json:"dropped,omitempty"`
}
