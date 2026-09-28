// Package rules evaluates Rule predicates deterministically over a fact map.
// Rules are executable constraints — the application enforces them; the
// model never gets to interpret them.
//
// Predicate grammar (v1):
//
//	path.to.field                 → dotted lookup into the fact map
//	count(path.to.list)           → length of a slice/map
//	a.b != null                   → existence checks
//	a.b == <number|string|path>   → equality
//	a.b >= <= > < <number|path>   → ordered comparison
//	a.b contains "x"              → substring / slice membership
package rules

import (
	"encoding/json"
	"fmt"
	"strings"

	"dramastudio/internal/intelligence/domain"
)

// Engine evaluates rules against task facts.
type Engine struct{}

func NewEngine() *Engine { return &Engine{} }

// Check runs every rule whose trigger matches (action, phase) against facts
// and returns violations. A nil rule list or empty Requires means "passes".
func (e *Engine) Check(rules []*domain.Rule, action, phase string, facts map[string]interface{}) []domain.RuleViolation {
	var out []domain.RuleViolation
	for _, r := range rules {
		if r == nil || !matchesTrigger(r.When, action, phase) {
			continue
		}
		for _, pred := range r.Requires {
			ok, err := evalPredicate(pred, facts)
			if err != nil {
				// A predicate that cannot evaluate is treated as a violation
				// evidence gap — reported rather than silently passed.
				out = append(out, violation(r, pred, fmt.Sprintf("predicate error: %v", err)))
				continue
			}
			if !ok {
				msg := r.OnFailure.Message
				if msg == "" {
					msg = fmt.Sprintf("rule %s violated: %s", r.ID, pred)
				}
				out = append(out, violation(r, pred, msg))
			}
		}
	}
	return out
}

func violation(r *domain.Rule, pred, msg string) domain.RuleViolation {
	return domain.RuleViolation{
		RuleID: r.ID, RuleVersion: r.Version, Predicate: pred,
		Severity: r.Severity, Enforcement: r.Enforcement,
		Message: msg, Phase: string(r.When.Phase),
	}
}

func matchesTrigger(t domain.RuleTrigger, action, phase string) bool {
	if t.Action != "" && t.Action != "*" && t.Action != action {
		return false
	}
	if t.Phase == "" {
		return true
	}
	return t.Phase == phase
}

// evalPredicate evaluates "lhs OP rhs". Unknown paths evaluate to nil.
func evalPredicate(pred string, facts map[string]interface{}) (bool, error) {
	pred = strings.TrimSpace(pred)
	if pred == "" {
		return true, nil
	}
	for _, op := range []string{"!=", "==", ">=", "<=", ">", "<", " contains "} {
		if i := strings.Index(pred, op); i > 0 {
			lhs := strings.TrimSpace(pred[:i])
			rhs := strings.TrimSpace(pred[i+len(op):])
			return compare(resolve(lhs, facts), resolve(rhs, facts), strings.TrimSpace(op))
		}
	}
	// Bare path → truthiness (non-nil, non-false, non-empty).
	v := resolve(pred, facts)
	return truthy(v), nil
}

// resolve evaluates a term: quoted string, number, boolean, null,
// count(path), or a dotted path lookup.
func resolve(term string, facts map[string]interface{}) interface{} {
	term = strings.TrimSpace(term)
	switch {
	case term == "null":
		return nil
	case term == "true":
		return true
	case term == "false":
		return false
	case strings.HasPrefix(term, "\"") && strings.HasSuffix(term, "\""):
		return strings.Trim(term, "\"")
	case strings.HasPrefix(term, "count(") && strings.HasSuffix(term, ")"):
		inner := term[6 : len(term)-1]
		return float64(countOf(resolve(inner, facts)))
	}
	if n, ok := parseNumber(term); ok {
		return n
	}
	return lookup(facts, term)
}

// lookup walks a dotted path through nested maps.
func lookup(facts map[string]interface{}, path string) interface{} {
	var cur interface{} = facts
	for _, seg := range strings.Split(path, ".") {
		m, ok := cur.(map[string]interface{})
		if !ok {
			return nil
		}
		cur = m[seg]
	}
	return cur
}

func compare(lhs, rhs interface{}, op string) (bool, error) {
	switch op {
	case "==":
		return equal(lhs, rhs), nil
	case "!=":
		return !equal(lhs, rhs), nil
	case "contains":
		return contains(lhs, rhs), nil
	}
	lf, lok := asFloat(lhs)
	rf, rok := asFloat(rhs)
	if !lok || !rok {
		return false, fmt.Errorf("non-numeric comparison: %v %s %v", lhs, op, rhs)
	}
	switch op {
	case ">=":
		return lf >= rf, nil
	case "<=":
		return lf <= rf, nil
	case ">":
		return lf > rf, nil
	case "<":
		return lf < rf, nil
	}
	return false, fmt.Errorf("unknown operator %q", op)
}

func equal(a, b interface{}) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	if af, ok := asFloat(a); ok {
		bf, bok := asFloat(b)
		return bok && af == bf
	}
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

func contains(lhs, rhs interface{}) bool {
	switch l := lhs.(type) {
	case string:
		s, _ := rhs.(string)
		return strings.Contains(l, s)
	case []interface{}:
		for _, item := range l {
			if equal(item, rhs) {
				return true
			}
		}
	}
	return false
}

func truthy(v interface{}) bool {
	switch t := v.(type) {
	case nil:
		return false
	case bool:
		return t
	case string:
		return t != ""
	case float64:
		return t != 0
	}
	return true
}

func countOf(v interface{}) int {
	switch c := v.(type) {
	case []interface{}:
		return len(c)
	case map[string]interface{}:
		return len(c)
	case string:
		return len(c)
	}
	return 0
}

func asFloat(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case json.Number:
		f, err := n.Float64()
		return f, err == nil
	}
	return 0, false
}

func parseNumber(s string) (float64, bool) {
	var f float64
	if _, err := fmt.Sscanf(s, "%g", &f); err == nil {
		// ensure full token consumed
		if fmt.Sprintf("%v", f) == s || strings.ContainsAny(s, ".eE") || isDigits(s) {
			return f, true
		}
	}
	return 0, false
}

func isDigits(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}
