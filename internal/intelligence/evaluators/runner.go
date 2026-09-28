// Package evaluators runs independent output assessment after generation —
// deterministic rule checks plus optional model-backed review — so the
// producer never grades its own work.
package evaluators

import (
	"context"
	"encoding/json"

	"dramastudio/internal/intelligence/domain"
	"dramastudio/internal/intelligence/rules"
)

// ModelReviewer is the seam for evaluators that need a model call (e.g.
// quality_evaluation capability). Injected so this package stays provider-
// agnostic and testable.
type ModelReviewer func(ctx context.Context, capability, prompt string) (string, error)

// Runner executes evaluators against generation output facts.
type Runner struct {
	rules  *rules.Engine
	cat    RuleLookup
	review ModelReviewer // may be nil — model evaluators are skipped
}

// RuleLookup resolves a rule ID to its definition (registry supplies it).
type RuleLookup func(id string) *domain.Rule

func NewRunner(lookup RuleLookup, review ModelReviewer) *Runner {
	return &Runner{rules: rules.NewEngine(), cat: lookup, review: review}
}

// Run executes one evaluator against outputFacts (the generated output plus
// the task fact bag). Rule checks are deterministic; a model capability adds
// a review call whose response becomes an INFO/WARNING finding.
func (r *Runner) Run(ctx context.Context, ev *domain.Evaluator, action string, outputFacts map[string]interface{}) ([]domain.EvalFinding, error) {
	var findings []domain.EvalFinding
	checkRules := make([]*domain.Rule, 0, len(ev.RuleIDs))
	for _, rid := range ev.RuleIDs {
		if rl := r.cat(rid); rl != nil {
			checkRules = append(checkRules, rl)
		}
	}
	for _, v := range r.rules.Check(checkRules, action, "post", outputFacts) {
		findings = append(findings, domain.EvalFinding{
			EvaluatorID: ev.ID, Check: v.RuleID, Severity: v.Severity, Detail: v.Message,
		})
	}
	if ev.ModelCapability != "" && r.review != nil {
		prompt := reviewPrompt(ev, outputFacts)
		resp, err := r.review(ctx, ev.ModelCapability, prompt)
		sev := domain.SeverityInfo
		detail := resp
		if err != nil {
			sev = domain.SeverityWarning
			detail = "model evaluation unavailable: " + err.Error()
		}
		findings = append(findings, domain.EvalFinding{
			EvaluatorID: ev.ID, Check: "model_review", Severity: sev, Detail: detail,
		})
	}
	return findings, nil
}

func reviewPrompt(ev *domain.Evaluator, facts map[string]interface{}) string {
	p := "Evaluate the " + ev.Target + " against these criteria:\n"
	for _, c := range ev.Rubric {
		p += "- " + c + "\n"
	}
	p += "\nOutput facts: " + marshalFacts(facts) + "\nRespond with PASS or REVISE plus findings."
	return p
}

func marshalFacts(m map[string]interface{}) string {
	b, _ := json.Marshal(m)
	return string(b)
}
