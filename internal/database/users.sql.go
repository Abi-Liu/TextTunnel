// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: users.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, username, password, created_at, updated_at, api_key) 
VALUES ($1, $2, $3, timezone('utc', NOW()), timezone('utc', NOW()), encode(sha256(random()::text::bytea), 'hex'))
RETURNING id, username, password, created_at, updated_at, api_key
`

type CreateUserParams struct {
	ID       uuid.UUID
	Username string
	Password string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.ID, arg.Username, arg.Password)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ApiKey,
	)
	return i, err
}

const deleteUserById = `-- name: DeleteUserById :execrows
DELETE FROM users WHERE id = $1
`

func (q *Queries) DeleteUserById(ctx context.Context, id uuid.UUID) (int64, error) {
	result, err := q.db.ExecContext(ctx, deleteUserById, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const findUserByApiKey = `-- name: FindUserByApiKey :one
SELECT id, username, password, created_at, updated_at, api_key FROM users WHERE api_key = $1
`

func (q *Queries) FindUserByApiKey(ctx context.Context, apiKey string) (User, error) {
	row := q.db.QueryRowContext(ctx, findUserByApiKey, apiKey)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ApiKey,
	)
	return i, err
}

const findUserById = `-- name: FindUserById :one
SELECT id, username, password, created_at, updated_at, api_key FROM users WHERE id = $1
`

func (q *Queries) FindUserById(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, findUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ApiKey,
	)
	return i, err
}

const findUserByUsername = `-- name: FindUserByUsername :one
SELECT id, username, password, created_at, updated_at, api_key FROM users WHERE username = $1
`

func (q *Queries) FindUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, findUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ApiKey,
	)
	return i, err
}
