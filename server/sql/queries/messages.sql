-- name: CreateMessage :one
INSERT INTO messages (id, content, created_at, updated_at, sender_id, room_id)
VALUES ($1, $2, timezone('utc', NOW()), timezone('utc', NOW()), $3, $4)
RETURNING *;

-- name: GetMessagesByUser :many
SELECT * from messages
WHERE sender_id = $1;

-- name: GetMessagesByRoom :many
SELECT * FROM messages WHERE room_id = $1;

-- name: GetMessagesByRoomAndUser :many
SELECT * FROM messages WHERE room_id = $1 AND sender_id = $2;
