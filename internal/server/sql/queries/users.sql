-- name: CreateUser :one
INSERT INTO users (id, username, password, created_at, updated_at) 
VALUES ($1, $2, $3, timezone('utc', NOW()), timezone('utc', NOW()))
RETURNING *;

-- name: FindUserById :one
SELECT * FROM users WHERE id = $1;

-- name: FindUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: DeleteUserById :execrows
DELETE FROM users WHERE id = $1;
