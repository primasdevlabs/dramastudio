package persistence

import (
	"context"
	"sync"
	"time"

	"dramastudio/internal/identity/domain"
)

type resetToken struct {
	userID    string
	expiresAt time.Time
	used      bool
}

type InMemoryUserRepository struct {
	mu          sync.RWMutex
	users       map[string]*domain.User
	byEmail     map[string]string
	orgs        map[string]*domain.Organization
	credentials map[string]*domain.APICredential
	byDigest    map[string]string
	memberships map[string]string      // "userID|orgID" -> role
	resets      map[string]*resetToken // tokenHash -> token
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users:       make(map[string]*domain.User),
		byEmail:     make(map[string]string),
		orgs:        make(map[string]*domain.Organization),
		credentials: make(map[string]*domain.APICredential),
		byDigest:    make(map[string]string),
		memberships: make(map[string]string),
		resets:      make(map[string]*resetToken),
	}
}

func (r *InMemoryUserRepository) SaveUser(_ context.Context, u *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if existing, ok := r.byEmail[u.Email]; ok && existing != u.ID {
		return domain.ErrEmailTaken
	}
	r.users[u.ID] = u
	r.byEmail[u.Email] = u.ID
	return nil
}

// withRole returns a copy of u with the membership role resolved — mirrors
// the SQL adapter's LEFT JOIN so callers see identical data on both drivers.
func (r *InMemoryUserRepository) withRole(u *domain.User) *domain.User {
	if role, ok := r.memberships[u.ID+"|"+u.OrgID]; ok {
		cpy := *u
		cpy.Role = role
		return &cpy
	}
	return u
}

func (r *InMemoryUserRepository) FindUserByID(_ context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.users[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	return r.withRole(u), nil
}

func (r *InMemoryUserRepository) FindUserByEmail(_ context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.byEmail[email]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	return r.withRole(r.users[id]), nil
}

func (r *InMemoryUserRepository) ListUsers(_ context.Context, orgID string) ([]*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.User, 0)
	for _, u := range r.users {
		if orgID == "" || u.OrgID == orgID {
			res = append(res, r.withRole(u))
		}
	}
	return res, nil
}

func (r *InMemoryUserRepository) SaveOrganization(_ context.Context, o *domain.Organization) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orgs[o.ID] = o
	return nil
}

func (r *InMemoryUserRepository) FindOrganizationByID(_ context.Context, id string) (*domain.Organization, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	o, ok := r.orgs[id]
	if !ok {
		return nil, domain.ErrOrgNotFound
	}
	return o, nil
}

func (r *InMemoryUserRepository) SaveMembership(_ context.Context, m *domain.Membership) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.memberships[m.UserID+"|"+m.OrganizationID] = m.Role
	return nil
}

func (r *InMemoryUserRepository) FindMembershipRole(_ context.Context, userID, orgID string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.memberships[userID+"|"+orgID], nil
}

func (r *InMemoryUserRepository) SaveCredential(_ context.Context, c *domain.APICredential) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.credentials[c.ID] = c
	r.byDigest[c.KeyDigest] = c.ID
	return nil
}

func (r *InMemoryUserRepository) FindCredentialByDigest(_ context.Context, digest string) (*domain.APICredential, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.byDigest[digest]
	if !ok {
		return nil, domain.ErrCredentialNotFound
	}
	c := r.credentials[id]
	if c.Revoked() {
		return nil, domain.ErrCredentialRevoked
	}
	return c, nil
}

func (r *InMemoryUserRepository) ListCredentials(_ context.Context, orgID string) ([]*domain.APICredential, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.APICredential, 0)
	for _, c := range r.credentials {
		if c.OrgID == orgID {
			res = append(res, c)
		}
	}
	return res, nil
}

func (r *InMemoryUserRepository) RevokeCredential(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.credentials[id]
	if !ok {
		return domain.ErrCredentialNotFound
	}
	now := time.Now().UTC()
	c.RevokedAt = &now
	return nil
}

func (r *InMemoryUserRepository) SaveResetToken(_ context.Context, userID, tokenHash string, expiresAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.resets[tokenHash]; !ok {
		r.resets[tokenHash] = &resetToken{userID: userID, expiresAt: expiresAt}
	}
	return nil
}

func (r *InMemoryUserRepository) FindUserByResetToken(_ context.Context, tokenHash string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.resets[tokenHash]
	if !ok || t.used || time.Now().UTC().After(t.expiresAt) {
		return nil, domain.ErrUserNotFound
	}
	u, ok := r.users[t.userID]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	return r.withRole(u), nil
}

func (r *InMemoryUserRepository) ConsumeResetToken(_ context.Context, tokenHash string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if t, ok := r.resets[tokenHash]; ok {
		t.used = true
	}
	return nil
}
