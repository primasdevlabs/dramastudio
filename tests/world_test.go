package tests

import (
	"context"
	"testing"

	worldsvc "dramastudio/internal/world/application/services"
	"dramastudio/internal/world/domain"
	worldinfra "dramastudio/internal/world/infrastructure/persistence"
)

func TestWorldLocationsVariantsProps(t *testing.T) {
	svc := worldsvc.NewWorldService(worldinfra.NewInMemoryWorldRepository())
	ctx := context.Background()

	parent, err := svc.CreateLocation(ctx, "proj_1", "House", "building", "", "")
	if err != nil {
		t.Fatalf("location: %v", err)
	}
	// Parent from another project is rejected.
	if _, err := svc.CreateLocation(ctx, "proj_1", "Room", "room", "", parent.ID); err != nil {
		t.Fatalf("same-project parent should pass: %v", err)
	}
	foreign, _ := svc.CreateLocation(ctx, "proj_x", "Elsewhere", "", "", "")
	if _, err := svc.CreateLocation(ctx, "proj_1", "Room2", "room", "", foreign.ID); err != domain.ErrLocationNotFound {
		t.Fatalf("foreign parent must be rejected, got %v", err)
	}

	if _, err := svc.AddVariant(ctx, parent.ID, "night", map[string]string{"lighting": "dim"}); err != nil {
		t.Fatalf("variant: %v", err)
	}
	got, _ := svc.GetLocation(ctx, parent.ID)
	if len(got.Variants) != 1 || got.Variants[0].Name != "night" {
		t.Fatalf("variant not hydrated: %+v", got.Variants)
	}

	if _, err := svc.CreateProp(ctx, "proj_1", "Knife", "", foreign.ID); err != domain.ErrLocationNotFound {
		t.Fatalf("foreign location prop must be rejected, got %v", err)
	}
	if _, err := svc.CreateProp(ctx, "proj_1", "Knife", "", parent.ID); err != nil {
		t.Fatalf("prop: %v", err)
	}
}
