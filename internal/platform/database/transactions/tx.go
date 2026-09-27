package transactions

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"dramastudio/internal/platform/database/postgres"
)

// TxManager runs a function inside a database transaction, committing on
// success and rolling back on error or panic.
type TxManager interface {
	WithTx(ctx context.Context, fn func(tx postgres.Querier) error) error
}

type pgxTxManager struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) TxManager {
	return &pgxTxManager{pool: pool}
}

func (m *pgxTxManager) WithTx(ctx context.Context, fn func(tx postgres.Querier) error) error {
	tx, err := m.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		// Rollback is a no-op if Commit already succeeded.
		_ = tx.Rollback(ctx)
	}()
	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
