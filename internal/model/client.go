package model

import "time"

type Client struct {
	ID               int64
	FirstName        string
	LastName         string
	Email            string
	ResidenceAddress string
	BirthDate        time.Time
	CreatedAt        time.Time
}
