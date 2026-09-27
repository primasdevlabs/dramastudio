package security

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService issues and validates HS256 access tokens carrying a Principal.
type JWTService struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTService(secret string, ttl time.Duration) (*JWTService, error) {
	if secret == "" {
		return nil, fmt.Errorf("security: JWT secret required")
	}
	if ttl <= 0 {
		return nil, fmt.Errorf("security: token TTL must be positive")
	}
	return &JWTService{secret: []byte(secret), ttl: ttl}, nil
}

type claims struct {
	OrgID       string   `json:"org_id"`
	Roles       []string `json:"roles,omitempty"`
	Permissions []string `json:"perms,omitempty"`
	Service     bool     `json:"svc,omitempty"`
	jwt.RegisteredClaims
}

func (s *JWTService) Issue(p Principal) (string, error) {
	c := claims{
		OrgID:       p.OrgID,
		Roles:       p.Roles,
		Permissions: p.Permissions,
		Service:     p.Service,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   p.UserID,
			Issuer:    "dramastudio",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.ttl)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(s.secret)
}

func (s *JWTService) ValidateToken(token string) (string, error) {
	p, err := s.Validate(token)
	if err != nil {
		return "", err
	}
	return p.UserID, nil
}

// Validate parses a token into a Principal.
func (s *JWTService) Validate(token string) (Principal, error) {
	var c claims
	_, err := jwt.ParseWithClaims(token, &c, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return s.secret, nil
	})
	if err != nil {
		return Principal{}, fmt.Errorf("invalid token: %w", err)
	}
	return Principal{
		UserID:      c.Subject,
		OrgID:       c.OrgID,
		Roles:       c.Roles,
		Permissions: c.Permissions,
		Service:     c.Service,
	}, nil
}
