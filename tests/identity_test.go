package tests

import (
	"context"
	"testing"

	idsvc "dramastudio/internal/identity/application/services"
	"dramastudio/internal/identity/domain"
	idinfra "dramastudio/internal/identity/infrastructure/persistence"
)

func newIdentityService(t *testing.T) *idsvc.IdentityService {
	t.Helper()
	return idsvc.NewIdentityService(idinfra.NewInMemoryUserRepository())
}

func TestIdentityRegisterPersistsOwnerMembership(t *testing.T) {
	svc := newIdentityService(t)
	ctx := context.Background()

	u, org, err := svc.Register(ctx, "Studio A", "owner@a.com", "Owner", "password123")
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if u.Role != domain.RoleOwner {
		t.Fatalf("expected owner role at register, got %q", u.Role)
	}

	// The role must survive a fresh read — regression test for the missing
	// identity.memberships write that silently downgraded owners to viewer.
	reloaded, err := svc.GetUser(ctx, org.ID, u.ID)
	if err != nil {
		t.Fatalf("get user: %v", err)
	}
	if reloaded.Role != domain.RoleOwner {
		t.Fatalf("membership role not persisted; got %q", reloaded.Role)
	}
}

func TestIdentityDuplicateEmailRejected(t *testing.T) {
	svc := newIdentityService(t)
	ctx := context.Background()

	if _, _, err := svc.Register(ctx, "Studio A", "dup@a.com", "One", "password123"); err != nil {
		t.Fatalf("first register: %v", err)
	}
	if _, _, err := svc.Register(ctx, "Studio B", "dup@a.com", "Two", "password123"); err != domain.ErrEmailTaken {
		t.Fatalf("expected ErrEmailTaken, got %v", err)
	}
}

func TestIdentityAuthenticate(t *testing.T) {
	svc := newIdentityService(t)
	ctx := context.Background()

	if _, _, err := svc.Register(ctx, "Studio A", "login@a.com", "User", "password123"); err != nil {
		t.Fatalf("register: %v", err)
	}
	u, err := svc.Authenticate(ctx, "login@a.com", "password123")
	if err != nil {
		t.Fatalf("authenticate: %v", err)
	}
	if u.Role != domain.RoleOwner {
		t.Fatalf("role lost on login read, got %q", u.Role)
	}
	if _, err := svc.Authenticate(ctx, "login@a.com", "wrong-password"); err != domain.ErrInvalidCredentials {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestIdentityAPIKeyLifecycle(t *testing.T) {
	svc := newIdentityService(t)
	ctx := context.Background()

	u, org, err := svc.Register(ctx, "Studio A", "keys@a.com", "User", "password123")
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	cred, plaintext, err := svc.CreateAPIKey(ctx, org.ID, u.ID, "ci", []string{"media.read"}, true)
	if err != nil {
		t.Fatalf("create key: %v", err)
	}
	if plaintext == "" || cred.KeyDigest == "" || cred.KeyDigest == plaintext {
		t.Fatal("plaintext must differ from stored digest")
	}
	resolved, err := svc.AuthenticateAPIKey(ctx, plaintext)
	if err != nil {
		t.Fatalf("authenticate key: %v", err)
	}
	if resolved.ID != cred.ID || resolved.OrgID != org.ID {
		t.Fatalf("wrong credential resolved: %+v", resolved)
	}
	if err := svc.RevokeCredential(ctx, org.ID, cred.ID); err != nil {
		t.Fatalf("revoke: %v", err)
	}
}

func TestIdentityCrossOrgAccessDenied(t *testing.T) {
	svc := newIdentityService(t)
	ctx := context.Background()

	uA, orgA, err := svc.Register(ctx, "Studio A", "a@a.com", "A", "password123")
	if err != nil {
		t.Fatalf("register A: %v", err)
	}
	_, orgB, err := svc.Register(ctx, "Studio B", "b@b.com", "B", "password123")
	if err != nil {
		t.Fatalf("register B: %v", err)
	}

	// Org B must not read org A's user or revoke org A's key.
	if _, err := svc.GetUser(ctx, orgB.ID, uA.ID); err != domain.ErrUserNotFound {
		t.Fatalf("cross-org getUser should hide the user, got %v", err)
	}
	_, _, err = svc.CreateAPIKey(ctx, orgA.ID, uA.ID, "a-key", nil, false)
	if err != nil {
		t.Fatalf("create key: %v", err)
	}
	creds, err := svc.ListCredentials(ctx, orgA.ID)
	if err != nil || len(creds) != 1 {
		t.Fatalf("list keys: %v len=%d", err, len(creds))
	}
	if err := svc.RevokeCredential(ctx, orgB.ID, creds[0].ID); err != domain.ErrUserNotFound {
		t.Fatalf("cross-org revoke should fail, got %v", err)
	}
}

func TestPasswordResetFlow(t *testing.T) {
	svc := newIdentityService(t)
	ctx := context.Background()

	if _, _, err := svc.Register(ctx, "Studio R", "reset@a.com", "Reset", "old-password"); err != nil {
		t.Fatalf("register: %v", err)
	}

	// Unknown email must not leak — empty token, nil error.
	if tok, err := svc.RequestPasswordReset(ctx, "nobody@a.com"); err != nil || tok != "" {
		t.Fatalf("unknown email should return empty token, got %q err=%v", tok, err)
	}

	token, err := svc.RequestPasswordReset(ctx, "reset@a.com")
	if err != nil || token == "" {
		t.Fatalf("expected reset token, got %q err=%v", token, err)
	}

	if err := svc.ResetPassword(ctx, token, "new-password-123"); err != nil {
		t.Fatalf("reset: %v", err)
	}
	if _, err := svc.Authenticate(ctx, "reset@a.com", "new-password-123"); err != nil {
		t.Fatalf("new password should authenticate: %v", err)
	}
	if _, err := svc.Authenticate(ctx, "reset@a.com", "old-password"); err == nil {
		t.Fatal("old password must no longer authenticate")
	}
	if err := svc.ResetPassword(ctx, token, "another-password"); err == nil {
		t.Fatal("token reuse must fail — tokens are single-use")
	}
}

func TestSSOAutoProvisioning(t *testing.T) {
	svc := newIdentityService(t)
	ctx := context.Background()

	u, created, err := svc.AuthenticateSSO(ctx, "sso@a.com", "SSO User", "")
	if err != nil || !created {
		t.Fatalf("first sign-in should provision, created=%v err=%v", created, err)
	}
	if u.Role != domain.RoleOwner {
		t.Fatalf("provisioned SSO user should be owner, got %q", u.Role)
	}
	u2, created, err := svc.AuthenticateSSO(ctx, "sso@a.com", "SSO User", "")
	if err != nil || created {
		t.Fatalf("second sign-in should reuse account, created=%v err=%v", created, err)
	}
	if u2.ID != u.ID {
		t.Fatal("SSO must resolve to the same user")
	}
}
