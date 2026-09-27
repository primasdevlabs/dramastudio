package configs

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Store backends supported by the runtime. "postgres" is the production
// default; "memory" exists for local development and tests only.
type StoreBackend string

const (
	StorePostgres StoreBackend = "postgres"
	StoreMemory   StoreBackend = "memory"
)

type StorageBackend string

const (
	StorageS3    StorageBackend = "s3"
	StorageLocal StorageBackend = "local"
)

type Config struct {
	Environment   string
	Server        ServerConfig
	Store         StoreBackend
	Database      DatabaseConfig
	Redis         RedisConfig
	Temporal      TemporalConfig
	Storage       StorageConfig
	Security      SecurityConfig
	Observability ObservabilityConfig
}

type ServerConfig struct {
	Port             int
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
	IdleTimeout      time.Duration
	ShutdownTimeout  time.Duration
	AllowedOrigins   []string
	PublicURL        string
}

type DatabaseConfig struct {
	URL          string
	MaxConns     int32
	MinConns     int32
	MigrationsDir string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type TemporalConfig struct {
	HostPort  string
	Namespace string
	APIKey    string
	UseTLS    bool
}

type StorageConfig struct {
	Backend   StorageBackend
	Endpoint  string
	Region    string
	Bucket    string
	AccessKey string
	SecretKey string
	UseSSL    bool
	LocalDir  string
}

type SecurityConfig struct {
	JWTSecret        string
	AccessTokenTTL   time.Duration
	ServiceTokenSalt string
	RequireAuth      bool
	DevUserID        string
	DevOrgID         string
}

type ObservabilityConfig struct {
	ServiceName   string
	TraceExporter string // "stdout" or "none"
	LogLevel      string
}

// Load reads configuration from the environment. Every value is explicit:
// missing required settings produce an error rather than a silent fallback.
func Load() (*Config, error) {
	cfg := &Config{
		Environment: envStr("DRAMASTUDIO_ENV", "development"),
		Server: ServerConfig{
			Port:            envInt("PORT", 8080),
			ReadTimeout:     envDur("HTTP_READ_TIMEOUT", 15*time.Second),
			WriteTimeout:    envDur("HTTP_WRITE_TIMEOUT", 60*time.Second),
			IdleTimeout:     envDur("HTTP_IDLE_TIMEOUT", 60*time.Second),
			ShutdownTimeout: envDur("HTTP_SHUTDOWN_TIMEOUT", 10*time.Second),
			AllowedOrigins:  envList("ALLOWED_ORIGINS"),
			PublicURL:       os.Getenv("PUBLIC_URL"),
		},
		Store: StoreBackend(envStr("DRAMASTUDIO_STORE", string(StorePostgres))),
		Database: DatabaseConfig{
			URL:           os.Getenv("DATABASE_URL"),
			MaxConns:      int32(envInt("DATABASE_MAX_CONNS", 10)),
			MinConns:      int32(envInt("DATABASE_MIN_CONNS", 1)),
			MigrationsDir: envStr("MIGRATIONS_DIR", "migrations"),
		},
		Redis: RedisConfig{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       envInt("REDIS_DB", 0),
		},
		Temporal: TemporalConfig{
			HostPort:  envStr("TEMPORAL_HOST_PORT", "localhost:7233"),
			Namespace: envStr("TEMPORAL_NAMESPACE", "default"),
			APIKey:    os.Getenv("TEMPORAL_API_KEY"),
			UseTLS:    envBool("TEMPORAL_TLS", false),
		},
		Storage: StorageConfig{
			Backend:   StorageBackend(envStr("STORAGE_BACKEND", string(StorageLocal))),
			Endpoint:  os.Getenv("S3_ENDPOINT"),
			Region:    envStr("S3_REGION", "us-east-1"),
			Bucket:    os.Getenv("S3_BUCKET"),
			AccessKey: os.Getenv("S3_ACCESS_KEY"),
			SecretKey: os.Getenv("S3_SECRET_KEY"),
			UseSSL:    envBool("S3_USE_SSL", true),
			LocalDir:  envStr("STORAGE_LOCAL_DIR", "data/objects"),
		},
		Security: SecurityConfig{
			JWTSecret:        os.Getenv("JWT_SECRET"),
			AccessTokenTTL:   envDur("ACCESS_TOKEN_TTL", 24*time.Hour),
			ServiceTokenSalt: os.Getenv("SERVICE_TOKEN_SALT"),
			RequireAuth:      envBool("REQUIRE_AUTH", false),
			DevUserID:        envStr("DEV_USER_ID", "dev-user"),
			DevOrgID:         envStr("DEV_ORG_ID", "dev-org"),
		},
		Observability: ObservabilityConfig{
			ServiceName:   envStr("OTEL_SERVICE_NAME", "dramastudio-api"),
			TraceExporter: envStr("OTEL_TRACE_EXPORTER", "none"),
			LogLevel:      envStr("LOG_LEVEL", "info"),
		},
	}
	return cfg, cfg.Validate()
}

func (c *Config) Validate() error {
	var missing []string

	switch c.Store {
	case StorePostgres:
		if c.Database.URL == "" {
			missing = append(missing, "DATABASE_URL")
		}
	case StoreMemory:
		// in-memory store needs no database
	default:
		return fmt.Errorf("invalid DRAMASTUDIO_STORE %q (want postgres|memory)", c.Store)
	}

	if c.Security.RequireAuth && c.Security.JWTSecret == "" {
		missing = append(missing, "JWT_SECRET")
	}

	switch c.Storage.Backend {
	case StorageS3:
		if c.Storage.Endpoint == "" {
			missing = append(missing, "S3_ENDPOINT")
		}
		if c.Storage.Bucket == "" {
			missing = append(missing, "S3_BUCKET")
		}
		if c.Storage.AccessKey == "" || c.Storage.SecretKey == "" {
			missing = append(missing, "S3_ACCESS_KEY/S3_SECRET_KEY")
		}
	case StorageLocal:
		if c.Storage.LocalDir == "" {
			missing = append(missing, "STORAGE_LOCAL_DIR")
		}
	default:
		return fmt.Errorf("invalid STORAGE_BACKEND %q (want s3|local)", c.Storage.Backend)
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required configuration: %s", strings.Join(missing, ", "))
	}
	return nil
}

func envStr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func envBool(key string, def bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return def
}

func envDur(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}

func envList(key string) []string {
	v := os.Getenv(key)
	if v == "" {
		return nil
	}
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if s := strings.TrimSpace(p); s != "" {
			out = append(out, s)
		}
	}
	return out
}
