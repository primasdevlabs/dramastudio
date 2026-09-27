package postgres

import (
	"context"
	"database/sql"
)

type DB struct {
	*sql.DB
}

func New(connStr string) (*DB, error) {
	return &DB{}, nil
}

func (db *DB) PingContext(ctx context.Context) error {
	return nil
}
