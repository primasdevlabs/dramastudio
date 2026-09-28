package tests

import (
	"context"
	"testing"

	charsvc "dramastudio/internal/characters/application/services"
	"dramastudio/internal/characters/domain"
	charinfra "dramastudio/internal/characters/infrastructure/persistence"
)

func newCharacterService(t *testing.T) *charsvc.CharacterService {
	t.Helper()
	return charsvc.NewCharacterService(charinfra.NewInMemoryCharacterRepository())
}

func TestCharacterVersionSnapshots(t *testing.T) {
	svc := newCharacterService(t)
	ctx := context.Background()

	c, err := svc.CreateCharacter(ctx, "proj_1", "Mara", "lead", "")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	// Save must snapshot v1 — same contract as the SQL adapter.
	vers, err := svc.ListVersions(ctx, c.ID)
	if err != nil || len(vers) != 1 {
		t.Fatalf("expected 1 version after create, got %d err=%v", len(vers), err)
	}
	if _, err := svc.SetAppearance(ctx, c.ID, domain.Appearance{DistinguishingMarks: "scarred"}); err != nil {
		t.Fatalf("appearance: %v", err)
	}
	vers, _ = svc.ListVersions(ctx, c.ID)
	if len(vers) != 2 {
		t.Fatalf("expected 2 versions, got %d", len(vers))
	}
	// v1 snapshot must not carry the v2 mutation.
	v1, err := svc.GetVersion(ctx, c.ID, 1)
	if err != nil {
		t.Fatalf("v1: %v", err)
	}
	if v1.Appearance.DistinguishingMarks == "scarred" {
		t.Fatal("v1 snapshot was mutated — versions must be immutable")
	}
}

func TestCharacterLockBlocksMutation(t *testing.T) {
	svc := newCharacterService(t)
	ctx := context.Background()

	c, _ := svc.CreateCharacter(ctx, "proj_1", "Mara", "lead", "")
	if err := svc.Lock(ctx, c.ID); err != nil {
		t.Fatalf("lock: %v", err)
	}
	if _, err := svc.SetAppearance(ctx, c.ID, domain.Appearance{DistinguishingMarks: "x"}); err != domain.ErrCharacterLocked {
		t.Fatalf("locked character must reject updates, got %v", err)
	}
	if err := svc.Unlock(ctx, c.ID); err != nil {
		t.Fatalf("unlock: %v", err)
	}
	if _, err := svc.SetAppearance(ctx, c.ID, domain.Appearance{DistinguishingMarks: "x"}); err != nil {
		t.Fatalf("unlocked update should work, got %v", err)
	}
}

func TestRelationshipRequiresSameProjectTarget(t *testing.T) {
	svc := newCharacterService(t)
	ctx := context.Background()

	a, _ := svc.CreateCharacter(ctx, "proj_1", "A", "lead", "")
	b, _ := svc.CreateCharacter(ctx, "proj_1", "B", "support", "")
	foreign, _ := svc.CreateCharacter(ctx, "proj_2", "C", "lead", "")

	if _, err := svc.AddRelationship(ctx, a.ID, domain.Relationship{
		TargetCharacterID: b.ID, RelationshipType: "sibling",
	}); err != nil {
		t.Fatalf("same-project relationship: %v", err)
	}
	if _, err := svc.AddRelationship(ctx, a.ID, domain.Relationship{
		TargetCharacterID: foreign.ID, RelationshipType: "rival",
	}); err != domain.ErrCharacterNotFound {
		t.Fatalf("cross-project relationship must be rejected, got %v", err)
	}
}
