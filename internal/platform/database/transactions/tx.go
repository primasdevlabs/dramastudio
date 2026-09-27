package transactions

import (
	"context"
	"database/sql"
)

type TxManager interface {
	WithTx(ctx context.Context, fn func(tx *sql.Tx) error) error
}
