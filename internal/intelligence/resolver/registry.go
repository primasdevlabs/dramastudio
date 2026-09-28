// Package resolver turns catalog definitions + DB policy layers into the
// AgentExecutionContext the executor consumes. It is the seam between the
// production pipeline and the intelligence layer.
package resolver

import (
	"fmt"
	"strings"

	"dramastudio/internal/intelligence/catalog"
	"dramastudio/internal/intelligence/domain"
)

// Registry resolves definitions by ID from the loaded catalog plus any
// operator-added policy layers in the repository.
type Registry struct {
	cat    *catalog.Catalog
	layers domain.PolicyLayerRepository
}

func NewRegistry(cat *catalog.Catalog, layers domain.PolicyLayerRepository) *Registry {
	return &Registry{cat: cat, layers: layers}
}

func (r *Registry) Agent(id string) (*domain.Agent, error) {
	if a, ok := r.cat.Agents[id]; ok {
		return a, nil
	}
	return nil, fmt.Errorf("agent %q not in catalog", id)
}

func (r *Registry) Agents() []*domain.Agent {
	out := make([]*domain.Agent, 0, len(r.cat.Agents))
	for _, a := range r.cat.Agents {
		out = append(out, a)
	}
	return out
}

func (r *Registry) Skills() []*domain.Skill {
	out := make([]*domain.Skill, 0, len(r.cat.Skills))
	for _, s := range r.cat.Skills {
		out = append(out, s)
	}
	return out
}

func (r *Registry) Rules() []*domain.Rule {
	out := make([]*domain.Rule, 0, len(r.cat.Rules))
	for _, rl := range r.cat.Rules {
		out = append(out, rl)
	}
	return out
}

func (r *Registry) Evaluators() []*domain.Evaluator {
	out := make([]*domain.Evaluator, 0, len(r.cat.Evaluators))
	for _, e := range r.cat.Evaluators {
		out = append(out, e)
	}
	return out
}

func (r *Registry) Policies() map[string][]*domain.Policy { return r.cat.Policies }

// ScopeIDs maps each policy layer to the concrete scope ID in effect for a
// task ("" for layers that don't apply). Built from task.Scope + ProjectID.
type ScopeIDs map[domain.PolicyLayer]string

// ScopeIDsFor parses "episode:ep-7" style scope strings into a layer map.
// The project layer always binds to task.ProjectID.
func ScopeIDsFor(task domain.TaskContext) ScopeIDs {
	ids := ScopeIDs{domain.LayerProject: task.ProjectID}
	if i := strings.Index(task.Scope, ":"); i > 0 {
		layer := domain.PolicyLayer(task.Scope[:i])
		ids[layer] = task.Scope[i+1:]
	}
	return ids
}
