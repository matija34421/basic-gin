package dto

import "time"

type AccountUpdate struct {
	ID            int      `json:"id" binding:"required"`
	ClientID      *int     `json:"client_id,omitempty"`
	AccountNumber *string  `json:"account_number,omitempty"`
	Balance       *float64 `json:"balance,omitempty"`
}

type AccountResponse struct {
	ID            int       `json:"id"`
	ClientID      int       `json:"client_id"`
	AccountNumber string    `json:"account_number"`
	Balance       float64   `json:"balance"`
	CreatedAt     time.Time `json:"created_at"`
}
