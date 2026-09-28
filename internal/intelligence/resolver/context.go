package resolver

import (
	"context"
	"fmt"

	"dramastudio/internal/intelligence/domain"
)

// ResolveExecution builds the AgentExecutionContext for a task: the agent
// definition, its resolved skills, effective policies (layered), applicable
// rules (explicit + any rule whose trigger matches the task action), and the
// evaluators guarding the output.
func (r *Registry) ResolveExecution(ctx context.Context, agentID string, task domain.TaskContext) (*domain.AgentExecutionContext, error) {
	agent, err := r.Agent(agentID)
	if err != nil {
		return nil, err
	}
	ec := &domain.AgentExecutionContext{
		Agent: agent, Task: task, Permissions: agent.Permissions,
	}
	scope := ScopeIDsFor(task)

	if err := r.collectSkills(ec, agent); err != nil {
		return nil, err
	}
	if err := r.collectPolicies(ctx, ec, agent, scope); err != nil {
		return nil, err
	}
	r.collectRules(ec, agent, task)
	if err := r.collectEvaluators(ec, agent); err != nil {
		return nil, err
	}
	return ec, nil
}

// collectSkills resolves the agent's skill IDs and seeds instruction lines
// with each skill's methodology.
func (r *Registry) collectSkills(ec *domain.AgentExecutionContext, agent *domain.Agent) error {
	for _, sid := range agent.SkillIDs {
		s, ok := r.cat.Skills[sid]
		if !ok {
			return fmt.Errorf("agent %s: skill %q not in catalog", agent.ID, sid)
		}
		ec.Skills = append(ec.Skills, s)
		ec.InstructionLines = append(ec.InstructionLines, s.Method...)
	}
	return nil
}

// collectPolicies merges every policy layer applicable to the task scope.
func (r *Registry) collectPolicies(ctx context.Context, ec *domain.AgentExecutionContext, agent *domain.Agent, scope ScopeIDs) error {
	for _, pid := range agent.PolicyIDs {
		eff, err := r.ResolvePolicy(ctx, pid, scope)
		if err != nil {
			return fmt.Errorf("agent %s: policy %s: %w", agent.ID, pid, err)
		}
		ec.Policies = append(ec.Policies, eff)
		ec.InstructionLines = append(ec.InstructionLines, eff.Constraints...)
	}
	return nil
}

// collectRules attaches explicit agent rules plus every catalog rule whose
// trigger targets this action.
func (r *Registry) collectRules(ec *domain.AgentExecutionContext, agent *domain.Agent, task domain.TaskContext) {
	seen := map[string]bool{}
	for _, rid := range agent.RuleIDs {
		if rl, ok := r.cat.Rules[rid]; ok {
			ec.Rules = append(ec.Rules, rl)
			seen[rid] = true
		}
	}
	for _, rl := range r.cat.Rules {
		if !seen[rl.ID] && (rl.When.Action == "*" || rl.When.Action == task.Action) {
			ec.Rules = append(ec.Rules, rl)
			seen[rl.ID] = true
		}
	}
}

// collectEvaluators resolves the evaluators guarding this agent's output.
func (r *Registry) collectEvaluators(ec *domain.AgentExecutionContext, agent *domain.Agent) error {
	for _, eid := range agent.EvaluatorIDs {
		ev, ok := r.cat.Evaluators[eid]
		if !ok {
			return fmt.Errorf("agent %s: evaluator %q not in catalog", agent.ID, eid)
		}
		ec.Evaluators = append(ec.Evaluators, ev)
	}
	return nil
}
