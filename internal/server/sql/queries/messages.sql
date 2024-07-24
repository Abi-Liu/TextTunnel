-- name: CreateMessage :one
INSERT INTO messages (id, content, created_at, updated_at, sender_id, room_id)
VALUES ($1, $2, timezone('utc', NOW()), timezone('utc', NOW()), $3, $4)
RETURNING *;

-- name: GetMessagesByUser :many
SELECT messages.*, users.username
FROM messages
JOIN users ON users.id = messages.sender_id
WHERE messages.sender_id = $1;

-- name: GetMessagesByRoom :many
SELECT messages.*, users.username 
FROM messages 
JOIN users ON users.id = messages.sender_id
WHERE messages.room_id = $1;

-- name: GetMessagesByRoomAndUser :many
SELECT * FROM messages WHERE room_id = $1 AND sender_id = $2;

-- name: GetPreviousRoomMessages :many
SELECT messages.*, users.username FROM messages 
JOIN users ON users.id = messages.sender_id
WHERE messages.room_id = $1 
ORDER BY messages.created_at ASC LIMIT $2;
