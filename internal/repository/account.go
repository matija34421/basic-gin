package repository

import (
	"basic-gin/internal/model"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountRepository struct {
	pool *pgxpool.Pool
}

func NewAccountRepository(pool *pgxpool.Pool) *AccountRepository {
	return &AccountRepository{
		pool: pool,
	}
}

func (r *AccountRepository) GetByClientId(ctx context.Context, id int) ([]*model.Account, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, client_id, account_number, balance, created_at FROM accounts WHERE client_id = $1", id)

	if err != nil {
		return nil, fmt.Errorf("get accounts by client id: %w", err)
	}

	defer rows.Close()

	var accounts []*model.Account

	for rows.Next() {
		var account model.Account

		if err := rows.Scan(&account.ID, &account.ClientId, &account.AccountNumber, &account.Balance, &account.CreatedAt); err != nil {
			return nil, fmt.Errorf("scanning rows: %v", err)
		}

		accounts = append(accounts, &account)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("account rows: %v", err)
	}

	return accounts, nil
}

func (r *AccountRepository) GetById(ctx context.Context, id int) (*model.Account, error) {
	var account model.Account

	if err := r.pool.QueryRow(ctx, "SELECT id, account_number, balance, client_id, created_at FROM accounts WHERE id = $1", id).Scan(&account.ID, &account.AccountNumber, &account.Balance, &account.ClientId, &account.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("account with an id: %d not found", id)
		}
		return nil, fmt.Errorf("get account by id: %v", err)
	}

	return &account, nil
}

func (r *AccountRepository) CreateAccount(ctx context.Context, account *model.Account) (*model.Account, error) {
	var savedAccount model.Account
	if err := r.pool.QueryRow(ctx, `INSERT INTO accounts(account_number, balance, client_id)
		values($1,$2,$3)
		RETURNING id, account_number, balance, client_id, created_at`,
		account.AccountNumber,
		account.Balance,
		account.ClientId,
	).Scan(
		&savedAccount.ID,
		&savedAccount.AccountNumber,
		&savedAccount.Balance,
		&savedAccount.ClientId,
		&savedAccount.CreatedAt,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, fmt.Errorf("account_number already exists")
		}
		return nil, fmt.Errorf("insert account: %w", err)
	}

	return &savedAccount, nil
}

func (r *AccountRepository) UpdateBalanceDelta(ctx context.Context, tx pgx.Tx, accountID int, delta float64) (*model.Account, error) {
	var acc model.Account
	err := tx.QueryRow(ctx, `
        UPDATE accounts
        SET balance = balance + $2
        WHERE id = $1 AND balance + $2 >= 0
        RETURNING id, account_number, balance, client_id, created_at
    `, accountID, delta).Scan(&acc.ID, &acc.AccountNumber, &acc.Balance, &acc.ClientId, &acc.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("insufficient funds or account not found")
		}
		return nil, fmt.Errorf("update balance: %w", err)
	}
	return &acc, nil
}

func (r *AccountRepository) Pool() *pgxpool.Pool { return r.pool }
