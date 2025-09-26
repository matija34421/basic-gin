package model

import "time"

type Transaction struct {
	ID            string    `db:"id"`
	FromAccountID int       `db:"from_account_id"`
	ToAccountID   int       `db:"to_account_id"`
	Amount        float64   `db:"amount"`
	CreatedAt     time.Time `db:"created_at"`
}
