-- name: CreateRoom :one
INSERT INTO rooms (id, name, created_at, updated_at, creator_id, owner_id)
VALUES ($1, $2, timezone('UTC', NOW()),timezone('UTC', NOW()), $3, $4)
RETURNING *;

-- name: FindAllRooms :many
SELECT * FROM rooms;

-- name: FindRoomById :one
SELECT * FROM rooms WHERE id = $1;

