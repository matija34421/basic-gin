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

type ClientRepository struct {
	pool *pgxpool.Pool
}

func NewClientRepository(pool *pgxpool.Pool) *ClientRepository {
	return &ClientRepository{
		pool: pool,
	}
}

func (r *ClientRepository) GetAll(ctx context.Context) ([]*model.Client, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, first_name, last_name, email, residence_address, birth_date, created_at FROM clients")

	if err != nil {
		return nil, fmt.Errorf("get all clients query: %v", err)
	}
	defer rows.Close()

	var clients []*model.Client

	for rows.Next() {
		var client model.Client

		if err := rows.Scan(&client.ID, &client.FirstName, &client.LastName, &client.Email, &client.ResidenceAddress, &client.BirthDate, &client.CreatedAt); err != nil {
			return nil, fmt.Errorf("scanning rows: %v", err)
		}

		clients = append(clients, &client)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("client rows: %v", err)
	}

	return clients, nil
}

func (r *ClientRepository) GetById(ctx context.Context, id int64) (*model.Client, error) {
	var c model.Client
	if err := r.pool.QueryRow(ctx, "SELECT id, first_name, last_name, email, residence_address, birth_date, created_at FROM clients WHERE id = $1", id).Scan(
		&c.ID,
		&c.FirstName,
		&c.LastName,
		&c.Email,
		&c.ResidenceAddress,
		&c.BirthDate,
		&c.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("client %d not found", id)
		}
		return nil, fmt.Errorf("get client by id query: %w", err)
	}

	return &c, nil
}

func (r *ClientRepository) CreateClient(ctx context.Context, client model.Client) (*model.Client, error) {
	var result model.Client
	err := r.pool.QueryRow(ctx, `INSERT INTO clients(first_name, last_name, email, residence_address, birth_date)
			values($1,$2,$3,$4,$5)
			RETURNING id, first_name, email, last_name, residence_address, birth_date`,
		client.FirstName,
		client.LastName,
		client.Email,
		client.ResidenceAddress,
		client.BirthDate,
	).Scan(&result.ID, &result.FirstName, &result.LastName, &result.Email, &result.ResidenceAddress, &result.BirthDate)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique email
			return nil, fmt.Errorf("email already exists")
		}
		return nil, fmt.Errorf("insert client: %w", err)
	}

	return &result, nil
}

func (r *ClientRepository) UpdateClient(ctx context.Context, client model.Client) (*model.Client, error) {
	var result model.Client
	err := r.pool.QueryRow(ctx, `UPDATE clients SET first_name = $1, last_name=$2, email=$3, residence_address=$4, birth_date=$5 WHERE id = $6 
			RETURNING id, first_name, last_name, email, residence_address, birth_date, created_at`,
		&client.FirstName,
		&client.LastName,
		&client.Email,
		&client.ResidenceAddress,
		&client.BirthDate,
		&client.CreatedAt,
	).Scan(&result.ID, &result.FirstName, &result.LastName, &client.Email, &result.ResidenceAddress, &result.BirthDate, &result.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("client %d not found", client.ID)
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			return nil, fmt.Errorf("email already exists")
		}
		return nil, fmt.Errorf("update client: %w", err)
	}
	return &result, nil
}

func (r *ClientRepository) DeleteClient(ctx context.Context, id int64) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM clients WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete client: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("client %d not found", id)
	}
	return nil
}
