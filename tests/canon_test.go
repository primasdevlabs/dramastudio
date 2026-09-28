package tests

import (
	"context"
	"testing"

	canonsvc "dramastudio/internal/canon/application/services"
	"dramastudio/internal/canon/domain"
	canoninfra "dramastudio/internal/canon/infrastructure/persistence"
)

func newCanonService(t *testing.T) *canonsvc.CanonService {
	t.Helper()
	return canonsvc.NewCanonService(canoninfra.NewInMemoryCanonRepository())
}

func TestCanonFactVersioning(t *testing.T) {
	svc := newCanonService(t)
	ctx := context.Background()

	f, err := svc.EstablishFact(ctx, &domain.StoryFact{
		ProjectID: "proj_1", Subject: "mara", Predicate: "knows", Object: "the_code",
	})
	if err != nil || f.Version != 1 || f.Status != domain.FactStatusCanonical {
		t.Fatalf("establish: %v %+v", err, f)
	}
	if _, err := svc.UpdateFact(ctx, f.ID, func(f *domain.StoryFact) error {
		f.Object = "the_new_code"
		return nil
	}); err != nil {
		t.Fatalf("update: %v", err)
	}
	vers, err := svc.ListFactVersions(ctx, "proj_1", f.ID)
	if err != nil || len(vers) != 2 {
		t.Fatalf("versions: %v len=%d", err, len(vers))
	}
	if vers[0].Version != 1 || vers[1].Version != 2 {
		t.Fatalf("version ordering wrong: %v", vers)
	}

	retconned, err := svc.Retcon(ctx, "proj_1", f.ID)
	if err != nil || retconned.Status != domain.FactStatusRetconned {
		t.Fatalf("retcon: %v status=%s", err, retconned.Status)
	}
	// Retcon appends history rather than rewriting it.
	vers, _ = svc.ListFactVersions(ctx, "proj_1", f.ID)
	if len(vers) != 3 {
		t.Fatalf("expected 3 versions after retcon, got %d", len(vers))
	}
}

func TestCanonCrossProjectFactDenied(t *testing.T) {
	svc := newCanonService(t)
	ctx := context.Background()

	f, _ := svc.EstablishFact(ctx, &domain.StoryFact{
		ProjectID: "proj_a", Subject: "s", Predicate: "p", Object: "o",
	})
	if _, err := svc.GetFact(ctx, "proj_b", f.ID); err != domain.ErrFactNotFound {
		t.Fatalf("cross-project getFact must 404, got %v", err)
	}
	if _, err := svc.Retcon(ctx, "proj_b", f.ID); err != domain.ErrFactNotFound {
		t.Fatalf("cross-project retcon must 404, got %v", err)
	}
}

func TestKnowledgeIsolationAndGrantScoping(t *testing.T) {
	svc := newCanonService(t)
	ctx := context.Background()

	f, _ := svc.EstablishFact(ctx, &domain.StoryFact{
		ProjectID: "proj_1", Subject: "s", Predicate: "p", Object: "o",
	})
	// Missing knowledge = knows nothing, not an error.
	ks, err := svc.GetKnowledge(ctx, "char_1", "ep_1")
	if err != nil || len(ks.KnownFactIDs) != 0 {
		t.Fatalf("empty knowledge: %v %+v", err, ks)
	}
	// Granting a fact from another project is rejected.
	foreign, _ := svc.EstablishFact(ctx, &domain.StoryFact{
		ProjectID: "proj_x", Subject: "s", Predicate: "p", Object: "o",
	})
	if _, err := svc.GrantKnowledge(ctx, "proj_1", "char_1", "ep_1", []string{foreign.ID}); err != domain.ErrFactNotFound {
		t.Fatalf("foreign fact grant must fail, got %v", err)
	}
	if _, err := svc.GrantKnowledge(ctx, "proj_1", "char_1", "ep_1", []string{f.ID}); err != nil {
		t.Fatalf("grant: %v", err)
	}
	knows, _ := svc.Knows(ctx, "char_1", "ep_1", f.ID)
	if !knows {
		t.Fatal("character should know the fact")
	}
	// Knowledge is per-episode — ep_2 must not inherit ep_1's knowledge.
	knows2, _ := svc.Knows(ctx, "char_1", "ep_2", f.ID)
	if knows2 {
		t.Fatal("knowledge leaked across episodes")
	}
}
