package httpjson

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	capability "dramastudio/internal/platform/ai/capability"
)

// Client is the shared HTTP plumbing for provider adapters. It normalizes
// transport failures and HTTP status codes into capability.ProviderError so
// every adapter produces the same error surface (Backend.md §84).
type Client struct {
	Provider     string
	BaseURL      string
	APIKey       string
	AuthStyle    AuthStyle
	ExtraHeaders map[string]string
	HTTP         *http.Client
}

type AuthStyle string

const (
	AuthBearer  AuthStyle = "bearer"
	AuthXAPIKey AuthStyle = "x-api-key"
	AuthNone    AuthStyle = "none"
)

func New(provider, baseURL, apiKey string, style AuthStyle, timeout time.Duration) *Client {
	return &Client{
		Provider:  provider,
		BaseURL:   strings.TrimRight(baseURL, "/"),
		APIKey:    apiKey,
		AuthStyle: style,
		HTTP:      &http.Client{Timeout: timeout},
	}
}

// Do sends a JSON request and decodes a JSON response.
func (c *Client) Do(ctx context.Context, method, path string, body, out interface{}) error {
	var rdr io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return &capability.ProviderError{Kind: capability.ErrInvalidRequest, Provider: c.Provider, Message: err.Error()}
		}
		rdr = bytes.NewReader(raw)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, rdr)
	if err != nil {
		return &capability.ProviderError{Kind: capability.ErrInvalidRequest, Provider: c.Provider, Message: err.Error()}
	}
	req.Header.Set("Content-Type", "application/json")
	c.applyAuth(req)
	for k, v := range c.ExtraHeaders {
		req.Header.Set(k, v)
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return mapTransportError(c.Provider, err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 32<<20))
	if err != nil {
		return &capability.ProviderError{Kind: capability.ErrUnknown, Provider: c.Provider, Message: err.Error(), Retryable: true}
	}
	if resp.StatusCode >= 400 {
		return mapHTTPError(c.Provider, resp, raw)
	}
	if out == nil {
		return nil
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return &capability.ProviderError{Kind: capability.ErrUnknown, Provider: c.Provider, Message: "decode response: " + err.Error()}
	}
	return nil
}

func (c *Client) applyAuth(req *http.Request) {
	switch c.AuthStyle {
	case AuthBearer:
		if c.APIKey != "" {
			req.Header.Set("Authorization", "Bearer "+c.APIKey)
		}
	case AuthXAPIKey:
		if c.APIKey != "" {
			req.Header.Set("x-api-key", c.APIKey)
		}
	}
}

func mapTransportError(provider string, err error) error {
	if isTimeout(err) {
		return &capability.ProviderError{Kind: capability.ErrTimeout, Provider: provider, Message: err.Error(), Retryable: true}
	}
	return &capability.ProviderError{Kind: capability.ErrUnavailable, Provider: provider, Message: err.Error(), Retryable: true}
}

func isTimeout(err error) bool {
	var ne net.Error
	return errors.As(err, &ne) && ne.Timeout()
}

func mapHTTPError(provider string, resp *http.Response, body []byte) error {
	msg := fmt.Sprintf("HTTP %d: %s", resp.StatusCode, truncate(string(body), 300))
	pe := &capability.ProviderError{Provider: provider, Message: msg}
	switch {
	case resp.StatusCode == 429:
		pe.Kind = capability.ErrRateLimited
		pe.Retryable = true
		if ra := resp.Header.Get("Retry-After"); ra != "" {
			if secs, err := strconv.Atoi(ra); err == nil {
				pe.RetryAfter = time.Duration(secs) * time.Second
			}
		}
	case resp.StatusCode == 408 || resp.StatusCode == 504:
		pe.Kind = capability.ErrTimeout
		pe.Retryable = true
	case resp.StatusCode >= 500:
		pe.Kind = capability.ErrUnavailable
		pe.Retryable = true
	case resp.StatusCode == 400 || resp.StatusCode == 422:
		pe.Kind = capability.ErrInvalidRequest
	case resp.StatusCode == 401 || resp.StatusCode == 403:
		pe.Kind = capability.ErrRejected
	default:
		pe.Kind = capability.ErrUnknown
	}
	return pe
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
