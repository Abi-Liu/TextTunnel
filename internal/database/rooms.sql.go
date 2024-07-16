// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: rooms.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createRoom = `-- name: CreateRoom :one
INSERT INTO rooms (id, name, created_at, updated_at, creator_id, owner_id)
VALUES ($1, $2, timezone('UTC', NOW()),timezone('UTC', NOW()), $3, $4)
RETURNING id, name, created_at, updated_at, creator_id, owner_id
`

type CreateRoomParams struct {
	ID        uuid.UUID
	Name      string
	CreatorID uuid.UUID
	OwnerID   uuid.UUID
}

func (q *Queries) CreateRoom(ctx context.Context, arg CreateRoomParams) (Room, error) {
	row := q.db.QueryRowContext(ctx, createRoom,
		arg.ID,
		arg.Name,
		arg.CreatorID,
		arg.OwnerID,
	)
	var i Room
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.CreatorID,
		&i.OwnerID,
	)
	return i, err
}

const findAllRooms = `-- name: FindAllRooms :many
SELECT id, name, created_at, updated_at, creator_id, owner_id FROM rooms
`

func (q *Queries) FindAllRooms(ctx context.Context) ([]Room, error) {
	rows, err := q.db.QueryContext(ctx, findAllRooms)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Room
	for rows.Next() {
		var i Room
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.CreatorID,
			&i.OwnerID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findRoomById = `-- name: FindRoomById :one
SELECT id, name, created_at, updated_at, creator_id, owner_id FROM rooms WHERE id = $1
`

func (q *Queries) FindRoomById(ctx context.Context, id uuid.UUID) (Room, error) {
	row := q.db.QueryRowContext(ctx, findRoomById, id)
	var i Room
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.CreatorID,
		&i.OwnerID,
	)
	return i, err
}
