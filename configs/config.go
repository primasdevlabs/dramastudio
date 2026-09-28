package configs

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Store backends supported by the runtime. "postgres" is the production
// database; "sqlite" is the local-development/test adapter (same repository
// code, translated dialect); "memory" is for tests only.
type StoreBackend string

const (
	StorePostgres StoreBackend = "postgres"
	StoreSQLite   StoreBackend = "sqlite"
	StoreMemory   StoreBackend = "memory"
)

type StorageBackend string

const (
	StorageS3    StorageBackend = "s3"
	StorageLocal StorageBackend = "local"
)

// WorkflowEngine selects the durable-execution adapter (§83). "temporal"
// requires a Temporal server; "local" runs the production pipeline
// in-process for fully-local development.
type WorkflowEngine string

const (
	EngineTemporal WorkflowEngine = "temporal"
	EngineLocal    WorkflowEngine = "local"
)

type Config struct {
	Environment   string
	Server        ServerConfig
	Store         StoreBackend
	Database      DatabaseConfig
	Redis         RedisConfig
	Workflow      WorkflowConfig
	Storage       StorageConfig
	Security      SecurityConfig
	Observability ObservabilityConfig
}

type ServerConfig struct {
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
	AllowedOrigins  []string
	PublicURL       string
}

type DatabaseConfig struct {
	URL           string // postgres DSN (DATABASE_URL)
	Path          string // sqlite file path (DATABASE_PATH)
	MaxConns      int32
	MinConns      int32
	MigrationsDir string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// WorkflowConfig picks the engine adapter plus engine-specific transport.
// Temporal fields apply only when Engine is "temporal".
type WorkflowConfig struct {
	Engine    WorkflowEngine
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
	GoogleClientID   string
}

type ObservabilityConfig struct {
	ServiceName   string
	TraceExporter string // "stdout" or "none"
	LogLevel      string
}

// Load reads configuration from the environment. Every value is explicit:
// missing required settings produce an error rather than a silent fallback.
//
// A .env file in the working directory or any parent (repo root when the
// binary runs from a subdirectory) is loaded first; variables already set in
// the real environment always win.
func Load() (*Config, error) {
	loadDotEnv()
	cfg := &Config{
		Environment: envStr("DRAMASTUDIO_ENV", "development"),
		Server: ServerConfig{
			Port:            envInt("PORT", 9471),
			ReadTimeout:     envDur("HTTP_READ_TIMEOUT", 15*time.Second),
			WriteTimeout:    envDur("HTTP_WRITE_TIMEOUT", 60*time.Second),
			IdleTimeout:     envDur("HTTP_IDLE_TIMEOUT", 60*time.Second),
			ShutdownTimeout: envDur("HTTP_SHUTDOWN_TIMEOUT", 10*time.Second),
			AllowedOrigins:  defaultOrigins(envList("ALLOWED_ORIGINS")),
			PublicURL:       os.Getenv("PUBLIC_URL"),
		},
		Store: StoreBackend(envStr("DRAMASTUDIO_STORE", defaultStore())),
		Database: DatabaseConfig{
			URL:           os.Getenv("DATABASE_URL"),
			Path:          envStr("DATABASE_PATH", "data/dramastudio.db"),
			MaxConns:      int32(envInt("DATABASE_MAX_CONNS", 10)),
			MinConns:      int32(envInt("DATABASE_MIN_CONNS", 1)),
			MigrationsDir: envStr("MIGRATIONS_DIR", "migrations"),
		},
		Redis: RedisConfig{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       envInt("REDIS_DB", 0),
		},
		Workflow: WorkflowConfig{
			Engine:    WorkflowEngine(envStr("WORKFLOW_ENGINE", defaultEngine())),
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
			GoogleClientID:   os.Getenv("GOOGLE_CLIENT_ID"),
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
	case StoreSQLite:
		// local/test profile — DATABASE_PATH defaults to data/dramastudio.db
	case StoreMemory:
		// in-memory store needs no database
	default:
		return fmt.Errorf("invalid DRAMASTUDIO_STORE %q (want postgres|sqlite|memory)", c.Store)
	}

	if c.Security.RequireAuth && c.Security.JWTSecret == "" {
		missing = append(missing, "JWT_SECRET")
	}

	switch c.Workflow.Engine {
	case EngineTemporal:
		if c.Workflow.HostPort == "" {
			missing = append(missing, "TEMPORAL_HOST_PORT")
		}
	case EngineLocal:
		// in-process pipeline needs no server
	default:
		return fmt.Errorf("invalid WORKFLOW_ENGINE %q (want temporal|local)", c.Workflow.Engine)
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

// defaultEngine picks the workflow engine that matches the deployment:
// development defaults to the in-process engine so a contributor needs
// only Go + a database file; production defaults to Temporal.
func defaultEngine() string {
	if envStr("DRAMASTUDIO_ENV", "development") == "development" {
		return string(EngineLocal)
	}
	return string(EngineTemporal)
}

// loadDotEnv populates missing environment variables from the nearest .env
// file, walking upward from the working directory. Existing env vars take
// precedence. Missing or malformed files are ignored — real env-only
// deployments are unaffected.
func loadDotEnv() {
	dir, err := os.Getwd()
	if err != nil {
		return
	}
	for {
		path := filepath.Join(dir, ".env")
		if f, err := os.Open(path); err == nil {
			defer f.Close()
			scan := bufio.NewScanner(f)
			for scan.Scan() {
				line := strings.TrimSpace(scan.Text())
				if line == "" || strings.HasPrefix(line, "#") {
					continue
				}
				key, value, ok := strings.Cut(line, "=")
				if !ok {
					continue
				}
				key = strings.TrimSpace(key)
				if key == "" || os.Getenv(key) != "" {
					continue // real environment wins
				}
				value = strings.TrimSpace(value)
				value = strings.Trim(value, `"'`)
				os.Setenv(key, value)
			}
			_ = scan.Err() // malformed/unreadable .env is non-fatal
			return
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return
		}
		dir = parent
	}
}

// defaultStore picks the persistence driver per deployment profile:
// development/test default to SQLite (self-contained, no server);
// production defaults to PostgreSQL.
func defaultStore() string {
	switch envStr("DRAMASTUDIO_ENV", "development") {
	case "production", "staging":
		return string(StorePostgres)
	default:
		return string(StoreSQLite)
	}
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

// defaultOrigins allows the studio dev server when running in development
// with no explicit allowlist. Production must set ALLOWED_ORIGINS.
func defaultOrigins(list []string) []string {
	if len(list) > 0 || envStr("DRAMASTUDIO_ENV", "development") != "development" {
		return list
	}
	return []string{"http://localhost:8742", "http://127.0.0.1:8742"}
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
