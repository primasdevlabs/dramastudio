package temporal

import (
	"context"
	"crypto/tls"
	"fmt"

	"go.temporal.io/sdk/client"
)

// Config mirrors configs.TemporalConfig without importing the root config
// package (keeps platform dependency-free).
type Config struct {
	HostPort  string
	Namespace string
	APIKey    string
	UseTLS    bool
}

// Dial connects to a Temporal server.
func Dial(_ context.Context, cfg Config) (client.Client, error) {
	opts := client.Options{
		HostPort:  cfg.HostPort,
		Namespace: cfg.Namespace,
	}
	if cfg.APIKey != "" {
		opts.Credentials = client.NewAPIKeyStaticCredentials(cfg.APIKey)
		opts.ConnectionOptions.TLS = &tls.Config{MinVersion: tls.VersionTLS12}
	} else if cfg.UseTLS {
		opts.ConnectionOptions.TLS = &tls.Config{MinVersion: tls.VersionTLS12}
	}
	c, err := client.Dial(client.Options(opts))
	if err != nil {
		return nil, fmt.Errorf("temporal: dial %s/%s: %w", cfg.HostPort, cfg.Namespace, err)
	}
	return c, nil
}
