package repository

import (
	"basic-gin/internal/model"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository struct {
	pool *pgxpool.Pool
}

func NewTransactionRepository(pool *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{pool: pool}
}

func (r *TransactionRepository) Pool() *pgxpool.Pool { return r.pool }

func (r *TransactionRepository) SaveTx(ctx context.Context, tx pgx.Tx, t *model.Transaction) error {
	return tx.QueryRow(ctx, `
		INSERT INTO transactions (from_account_id, to_account_id, amount)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`, t.FromAccountID, t.ToAccountID, t.Amount).
		Scan(&t.ID, &t.CreatedAt)
}

func (r *TransactionRepository) ListByAccountID(ctx context.Context, accountID, limit, offset int) ([]*model.Transaction, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, from_account_id, to_account_id, amount, created_at
		FROM transactions
		WHERE from_account_id = $1 OR to_account_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, accountID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*model.Transaction
	for rows.Next() {
		var t model.Transaction
		if err := rows.Scan(&t.ID, &t.FromAccountID, &t.ToAccountID, &t.Amount, &t.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, &t)
	}
	return out, rows.Err()
}
