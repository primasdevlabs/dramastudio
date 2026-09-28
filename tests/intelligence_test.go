package tests

import (
	"context"
	"encoding/json"
	"testing"

	"dramastudio/internal/intelligence/catalog"
	inteldomain "dramastudio/internal/intelligence/domain"
	"dramastudio/internal/intelligence/evaluators"
	"dramastudio/internal/intelligence/execution"
	intelpersistence "dramastudio/internal/intelligence/infrastructure/persistence"
	"dramastudio/internal/intelligence/resolver"
	"dramastudio/internal/intelligence/rules"
)

func loadCatalog(t *testing.T) *catalog.Catalog {
	t.Helper()
	cat, err := catalog.LoadDir("../catalogs/intelligence")
	if err != nil {
		t.Fatalf("load catalog: %v", err)
	}
	return cat
}

// TestCatalogLoads proves the shipped definition tree parses and covers the
// agent roster.
func TestCatalogLoads(t *testing.T) {
	cat := loadCatalog(t)
	if len(cat.Agents) != 19 {
		t.Fatalf("expected 19 agents, got %d", len(cat.Agents))
	}
	if len(cat.Skills) < 25 || len(cat.Rules) < 5 || len(cat.Evaluators) < 3 {
		t.Fatalf("thin catalog: skills=%d rules=%d evaluators=%d",
			len(cat.Skills), len(cat.Rules), len(cat.Evaluators))
	}
}

// TestRulesEngine covers deterministic predicate evaluation.
func TestRulesEngine(t *testing.T) {
	e := rules.NewEngine()
	rule := &inteldomain.Rule{
		ID: "r1", Version: 1,
		When:        inteldomain.RuleTrigger{Action: "generate_character_video"},
		Requires:    []string{"character.canonical_reference != null"},
		Severity:    inteldomain.SeverityBlocking,
		Enforcement: inteldomain.EnforceBlock,
		OnFailure:   inteldomain.RuleFailure{Message: "ref required"},
	}
	facts := map[string]interface{}{"character": map[string]interface{}{}}
	vs := e.Check([]*inteldomain.Rule{rule}, "generate_character_video", "pre", facts)
	if len(vs) != 1 || vs[0].Enforcement != inteldomain.EnforceBlock {
		t.Fatalf("expected block violation, got %+v", vs)
	}
	facts["character"] = map[string]interface{}{"canonical_reference": "ref-1"}
	if vs := e.Check([]*inteldomain.Rule{rule}, "generate_character_video", "pre", facts); len(vs) != 0 {
		t.Fatalf("expected pass, got %+v", vs)
	}
	// count() + numeric comparison
	rule2 := &inteldomain.Rule{
		ID: "r2", When: inteldomain.RuleTrigger{Action: "render_episode"},
		Requires:    []string{"count(continuity.issues) == 0"},
		Enforcement: inteldomain.EnforceBlock,
	}
	bad := map[string]interface{}{"continuity": map[string]interface{}{
		"issues": []interface{}{"a"}}}
	if vs := e.Check([]*inteldomain.Rule{rule2}, "render_episode", "pre", bad); len(vs) != 1 {
		t.Fatalf("expected count violation, got %+v", vs)
	}
}

type stubGen struct{ called bool }

func (g *stubGen) Generate(ctx context.Context, task inteldomain.TaskContext, prompt string) (*execution.GenResult, error) {
	g.called = true
	out, _ := json.Marshal(map[string]interface{}{"text": "ok"})
	return &execution.GenResult{
		JobID: "j1", ProviderID: "mock", ModelID: "mock.universal",
		ModelVersion: "1", Output: out, Status: "succeeded",
	}, nil
}

// TestExecutorBlocksOnRule proves a BLOCKING rule halts before generation.
func TestExecutorBlocksOnRule(t *testing.T) {
	cat := loadCatalog(t)
	mem := intelpersistence.NewInMemoryRepository()
	reg := resolver.NewRegistry(cat, mem)
	gen := &stubGen{}
	runner := evaluators.NewRunner(
		func(id string) *inteldomain.Rule { return cat.Rules[id] }, nil)
	exec := execution.NewExecutor(reg, runner, gen, mem)

	rec, res, err := exec.Run(context.Background(), "video-director", inteldomain.TaskContext{
		ProjectID: "p1", Action: "generate_character_video",
		Capability: "video_generation",
		Facts:      map[string]interface{}{}, // no character.canonical_reference
	})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if gen.called {
		t.Fatal("generator must not be called when a BLOCKING rule fails")
	}
	if res != nil || rec.Verdict != inteldomain.VerdictBlocked {
		t.Fatalf("expected blocked verdict, got %s res=%v", rec.Verdict, res)
	}
	stored, err := mem.FindExecution(context.Background(), rec.ID)
	if err != nil || stored.AgentVersion != 1 {
		t.Fatalf("audit record: %v %+v", err, stored)
	}
}

// TestExecutorPassesAndEvaluates exercises a clean run with evaluators.
func TestExecutorPassesAndEvaluates(t *testing.T) {
	cat := loadCatalog(t)
	mem := intelpersistence.NewInMemoryRepository()
	reg := resolver.NewRegistry(cat, mem)
	gen := &stubGen{}
	runner := evaluators.NewRunner(
		func(id string) *inteldomain.Rule { return cat.Rules[id] }, nil)
	exec := execution.NewExecutor(reg, runner, gen, mem)

	rec, res, err := exec.Run(context.Background(), "episode-writer", inteldomain.TaskContext{
		ProjectID: "p1", Action: "write_episode_script",
		Capability: "script_writing", Prompt: "Write cold open.",
		Facts: map[string]interface{}{},
	})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if !gen.called || res == nil || res.ProviderID != "mock" {
		t.Fatalf("generation path broken: %v", res)
	}
	if rec.Verdict != inteldomain.VerdictApproved && rec.Verdict != inteldomain.VerdictNeedsReview {
		t.Fatalf("unexpected verdict %s violations=%+v findings=%+v",
			rec.Verdict, rec.Violations, rec.Findings)
	}
	if rec.PromptHash == "" || len(rec.SkillIDs) == 0 {
		t.Fatalf("audit incomplete: %+v", rec)
	}
}

// TestPolicyLayering proves protected keys block lower-layer overrides.
func TestPolicyLayering(t *testing.T) {
	cat := loadCatalog(t)
	mem := intelpersistence.NewInMemoryRepository()
	reg := resolver.NewRegistry(cat, mem)
	ctx := context.Background()

	// Operator adds a project layer trying to weaken a protected key.
	_ = mem.SavePolicyLayer(ctx, &inteldomain.Policy{
		ID: "character-consistency", Version: 1, Layer: inteldomain.LayerProject,
		ScopeID: "p1",
		Requirements: map[string]interface{}{
			"canonical_reference_required": false, // protected at system layer
			"custom_project_rule":          true,
		},
	})
	eff, err := reg.ResolvePolicy(ctx, "character-consistency", resolver.ScopeIDs{
		inteldomain.LayerProject: "p1",
	})
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if eff.Requirements["canonical_reference_required"] != true {
		t.Fatalf("protected key was weakened: %+v", eff.Requirements)
	}
	if eff.Requirements["custom_project_rule"] != true {
		t.Fatalf("non-protected override lost: %+v", eff.Requirements)
	}
	if len(eff.Dropped) == 0 {
		t.Fatal("dropped override not recorded in trace")
	}
}
