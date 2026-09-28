// Package sqlite provides a development/test database adapter over
// modernc.org/sqlite (pure Go, no CGO). It implements postgres.Querier so
// every SQL repository runs unchanged; translate() rewrites the small
// Postgres dialect surface ($N placeholders, schema-qualified names, casts,
// JSONB containment) into SQLite at the boundary.
//
// Scope: local development, repository tests, and basic integration tests.
// It is not a production database — production runs on PostgreSQL.
package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	_ "modernc.org/sqlite"

	"dramastudio/internal/platform/database/postgres"
)

// DB wraps a database/sql handle on the modernc sqlite driver.
type DB struct {
	db *sql.DB
}

// Open opens (creating if needed) the database at path. Pass ":memory:" or a
// file path; parent directories must exist.
func Open(ctx context.Context, path string) (*DB, error) {
	if path == "" {
		return nil, fmt.Errorf("sqlite: empty path")
	}
	dsn := path
	if path != ":memory:" && !strings.HasPrefix(path, "file:") {
		abs, err := filepath.Abs(path)
		if err != nil {
			return nil, fmt.Errorf("sqlite: resolve path: %w", err)
		}
		dsn = "file:" + filepath.ToSlash(abs)
	}
	sep := "?"
	if strings.Contains(dsn, "?") {
		sep = "&"
	}
	dsn += sep + url.Values{
		"_pragma": {"foreign_keys(1)", "busy_timeout(5000)", "journal_mode(WAL)"},
	}.Encode()

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("sqlite: open: %w", err)
	}
	// A single connection serializes writers — WAL still allows concurrent
	// readers, and SQLite serializes writers internally anyway.
	db.SetMaxOpenConns(1)
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("sqlite: ping: %w", err)
	}
	return &DB{db: db}, nil
}

func (d *DB) PingContext(ctx context.Context) error { return d.db.PingContext(ctx) }
func (d *DB) Close()                                { d.db.Close() }

// SQL returns the underlying *sql.DB for migrations and direct checks.
func (d *DB) SQL() *sql.DB { return d.db }

// Exec implements postgres.Querier.
func (d *DB) Exec(ctx context.Context, sqlText string, args ...interface{}) (pgconn.CommandTag, error) {
	q, bound := translate(sqlText, args)
	res, err := d.db.ExecContext(ctx, q, bound...)
	if err != nil {
		return pgconn.CommandTag{}, err
	}
	n, _ := res.RowsAffected()
	return pgconn.NewCommandTag(fmt.Sprintf("EXEC %d", n)), nil
}

// Query implements postgres.Querier.
func (d *DB) Query(ctx context.Context, sqlText string, args ...interface{}) (pgx.Rows, error) {
	q, bound := translate(sqlText, args)
	rows, err := d.db.QueryContext(ctx, q, bound...)
	if err != nil {
		return nil, err
	}
	return &Rows{rows: rows}, nil
}

// QueryRow implements postgres.Querier.
func (d *DB) QueryRow(ctx context.Context, sqlText string, args ...interface{}) pgx.Row {
	q, bound := translate(sqlText, args)
	return &Row{row: d.db.QueryRowContext(ctx, q, bound...)}
}

var _ postgres.Querier = (*DB)(nil)

// Row adapts *sql.Row to pgx.Row, mapping sql.ErrNoRows to pgx.ErrNoRows so
// repositories' postgres.IsNoRows checks work unchanged.
type Row struct{ row *sql.Row }

func (r *Row) Scan(dest ...interface{}) error {
	err := r.row.Scan(dest...)
	if err != nil && err == sql.ErrNoRows {
		return pgx.ErrNoRows
	}
	return err
}

// Rows adapts *sql.Rows to pgx.Rows.
type Rows struct {
	rows     *sql.Rows
	consumed bool
}

func (r *Rows) Close()                                       { r.rows.Close() }
func (r *Rows) Err() error                                   { return r.rows.Err() }
func (r *Rows) CommandTag() pgconn.CommandTag                { return pgconn.CommandTag{} }
func (r *Rows) FieldDescriptions() []pgconn.FieldDescription { return nil }
func (r *Rows) Conn() *pgx.Conn                              { return nil }
func (r *Rows) RawValues() [][]byte                          { return nil }

func (r *Rows) Next() bool {
	if !r.rows.Next() {
		return false
	}
	return true
}

func (r *Rows) Scan(dest ...interface{}) error { return r.rows.Scan(dest...) }

// Values returns the current row's values using column names for sizing.
func (r *Rows) Values() ([]interface{}, error) {
	cols, err := r.rows.Columns()
	if err != nil {
		return nil, err
	}
	dest := make([]interface{}, len(cols))
	ptrs := make([]interface{}, len(cols))
	for i := range dest {
		ptrs[i] = &dest[i]
	}
	if err := r.rows.Scan(ptrs...); err != nil {
		return nil, err
	}
	return dest, nil
}

var _ pgx.Rows = (*Rows)(nil)
var _ pgx.Row = (*Row)(nil)

// Now formats a timestamp the way SQLite TIMESTAMP columns store it.
func Now() time.Time { return time.Now().UTC() }
