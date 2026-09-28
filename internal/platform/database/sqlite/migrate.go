package sqlite

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Migrate applies the PostgreSQL migration files in dir to the SQLite
// database, translating DDL on the way in. The same migrations/ directory
// stays the single schema source of truth for both drivers.
func Migrate(ctx context.Context, d *DB, dir string) error {
	if _, err := d.db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS schema_migrations (version TEXT PRIMARY KEY, applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP)`); err != nil {
		return fmt.Errorf("sqlite: migrations table: %w", err)
	}
	files, err := filepath.Glob(filepath.Join(dir, "*.sql"))
	if err != nil {
		return fmt.Errorf("sqlite: glob migrations: %w", err)
	}
	sort.Strings(files)
	for _, f := range files {
		version := strings.TrimSuffix(filepath.Base(f), ".sql")
		var applied int
		if err := d.db.QueryRowContext(ctx,
			`SELECT COUNT(*) FROM schema_migrations WHERE version = ?`, version).Scan(&applied); err != nil {
			return err
		}
		if applied > 0 {
			continue
		}
		data, err := os.ReadFile(f)
		if err != nil {
			return fmt.Errorf("sqlite: read %s: %w", version, err)
		}
		if err := applyMigration(ctx, d, version, string(data)); err != nil {
			return err
		}
	}
	return nil
}

func applyMigration(ctx context.Context, d *DB, version, body string) error {
	for _, stmt := range strings.Split(stripLineComments(body), ";") {
		stmt = translateDDL(stmt)
		if stmt == "" {
			continue
		}
		if _, err := d.db.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("sqlite: migration %s: %w\n  statement: %s", version, err, stmt)
		}
	}
	_, err := d.db.ExecContext(ctx, `INSERT INTO schema_migrations (version) VALUES (?)`, version)
	return err
}

var (
	ddlSchemaName = regexp.MustCompile(`\b(` + strings.Join(schemaNames, "|") + `)\.`)
	ddlNow        = regexp.MustCompile(`(?i)\bnow\(\s*\)`)
	ddlTypeMap    = map[*regexp.Regexp]string{
		regexp.MustCompile(`(?i)\bJSONB\b`):                                       `TEXT`,
		regexp.MustCompile(`(?i)\bTIMESTAMPTZ\b`):                                 `TIMESTAMP`,
		regexp.MustCompile(`(?i)\bTIMESTAMP\s+WITH\s+TIME\s+ZONE\b`):              `TIMESTAMP`,
		regexp.MustCompile(`(?i)\bBIGINT\s+GENERATED\s+ALWAYS\s+AS\s+IDENTITY\b`): `INTEGER`,
		regexp.MustCompile(`(?i)\bDOUBLE\s+PRECISION\b`):                          `REAL`,
		regexp.MustCompile(`(?i)\bSERIAL\b`):                                      `INTEGER`,
	}
)

// translateDDL converts one PostgreSQL DDL statement to SQLite, returning ""
// for statements SQLite cannot express (schema creation, GIN indexes).
func translateDDL(stmt string) string {
	upper := strings.ToUpper(stmt)
	switch {
	case strings.Contains(upper, "CREATE SCHEMA"),
		strings.Contains(upper, "CREATE EXTENSION"),
		strings.Contains(upper, "COMMENT ON"),
		strings.Contains(upper, "USING GIN"),
		strings.Contains(upper, "USING GIST"):
		return ""
	}
	if emptyStatement(stmt) {
		return ""
	}
	stmt = ddlSchemaName.ReplaceAllString(stmt, `${1}_`)
	for re, to := range ddlTypeMap {
		stmt = re.ReplaceAllString(stmt, to)
	}
	stmt = ddlNow.ReplaceAllString(stmt, `CURRENT_TIMESTAMP`)
	return strings.TrimSpace(stmt)
}

// stripLineComments removes -- comments before statement splitting so a
// semicolon inside a comment does not truncate a statement. Migration DDL
// contains no -- sequences inside string literals.
func stripLineComments(body string) string {
	var b strings.Builder
	for _, line := range strings.Split(body, "\n") {
		if i := strings.Index(line, "--"); i >= 0 {
			line = line[:i]
		}
		b.WriteString(line)
		b.WriteByte('\n')
	}
	return b.String()
}

// emptyStatement reports whether stmt has no executable content.
func emptyStatement(stmt string) bool {
	return strings.TrimSpace(stmt) == ""
}
