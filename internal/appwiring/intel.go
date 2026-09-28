package appwiring

import (
	"context"
	"encoding/json"

	intelapp "dramastudio/internal/intelligence/application"
	intelcatalog "dramastudio/internal/intelligence/catalog"
	inteldomain "dramastudio/internal/intelligence/domain"
	intelevals "dramastudio/internal/intelligence/evaluators"
	intelexec "dramastudio/internal/intelligence/execution"
	intelresolver "dramastudio/internal/intelligence/resolver"
)

// buildIntelligence wires the production-intelligence subsystem: catalog →
// registry → executor (wrapping the AI generation pipeline) → service.
// Definition files live at INTELLIGENCE_DIR (default catalogs/intelligence).
func (c *Container) buildIntelligence(r repoBundle) error {
	dir := envOr("INTELLIGENCE_DIR", "catalogs/intelligence")
	cat, err := intelcatalog.LoadDir(dir)
	if err != nil {
		return err
	}
	reg := intelresolver.NewRegistry(cat, r.intelLayers)
	gen := intelapp.NewAIGenerator(c.AIExecutor)
	review := modelReviewer(gen)
	runner := intelevals.NewRunner(
		func(id string) *inteldomain.Rule { return cat.Rules[id] }, review)
	exec := intelexec.NewExecutor(reg, runner, gen, r.intelExec)
	c.Intel = intelapp.NewIntelligenceService(reg, exec, r.intelExec, r.intelLayers)
	return nil
}

// modelReviewer lets evaluators issue a review call through the normal
// generation path (policy routing + provenance included).
func modelReviewer(gen intelexec.Generator) intelevals.ModelReviewer {
	return func(ctx context.Context, capability, prompt string) (string, error) {
		res, err := gen.Generate(ctx, inteldomain.TaskContext{
			Capability: capability, Prompt: prompt,
		}, prompt)
		if err != nil {
			return "", err
		}
		var out struct {
			Text string `json:"text"`
		}
		if err := json.Unmarshal(res.Output, &out); err != nil {
			return "", err
		}
		return out.Text, nil
	}
}
