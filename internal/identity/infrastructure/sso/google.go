// Package sso verifies third-party identity tokens server-side.
package sso

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

var ErrInvalidToken = errors.New("identity token failed verification")

// GoogleProfile is the verified identity extracted from a Google ID token.
type GoogleProfile struct {
	Subject       string
	Email         string
	Name          string
	EmailVerified bool
}

type tokenInfo struct {
	Issuer        string `json:"iss"`
	Audience      string `json:"aud"`
	Subject       string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Name          string `json:"name"`
	Expires       string `json:"exp"`
}

// GoogleVerifier validates Google ID tokens against Google's tokeninfo
// endpoint and pins the audience to the configured OAuth client ID.
type GoogleVerifier struct {
	ClientID string
	// Endpoint is overridable for tests.
	Endpoint   string
	HTTPClient *http.Client
}

func NewGoogleVerifier(clientID string) *GoogleVerifier {
	return &GoogleVerifier{
		ClientID:   clientID,
		Endpoint:   "https://oauth2.googleapis.com/tokeninfo",
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// Verify returns the token's identity claims or ErrInvalidToken.
func (v *GoogleVerifier) Verify(ctx context.Context, idToken string) (*GoogleProfile, error) {
	if v.ClientID == "" || idToken == "" {
		return nil, ErrInvalidToken
	}
	endpoint := v.Endpoint
	if endpoint == "" {
		endpoint = "https://oauth2.googleapis.com/tokeninfo"
	}
	hc := v.HTTPClient
	if hc == nil {
		hc = &http.Client{Timeout: 10 * time.Second}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		endpoint+"?id_token="+url.QueryEscape(idToken), nil)
	if err != nil {
		return nil, err
	}
	resp, err := hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token verification request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, ErrInvalidToken
	}

	var info tokenInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, ErrInvalidToken
	}

	switch info.Issuer {
	case "accounts.google.com", "https://accounts.google.com":
	default:
		return nil, ErrInvalidToken
	}
	if info.Audience != v.ClientID {
		return nil, ErrInvalidToken
	}
	if info.Email == "" || info.EmailVerified != "true" {
		return nil, ErrInvalidToken
	}
	if exp, err := parseUnixSeconds(info.Expires); err == nil && time.Now().After(exp) {
		return nil, ErrInvalidToken
	}

	return &GoogleProfile{
		Subject:       info.Subject,
		Email:         info.Email,
		Name:          info.Name,
		EmailVerified: true,
	}, nil
}

func parseUnixSeconds(s string) (time.Time, error) {
	var n int64
	if _, err := fmt.Sscanf(s, "%d", &n); err != nil {
		return time.Time{}, err
	}
	return time.Unix(n, 0), nil
}
