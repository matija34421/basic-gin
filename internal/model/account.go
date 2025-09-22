package model

import "time"

type Account struct {
	ID            int
	ClientId      int
	AccountNumber string
	Balance       float64
	CreatedAt     time.Time
}
