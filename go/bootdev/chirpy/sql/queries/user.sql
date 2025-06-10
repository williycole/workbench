-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email)
VALUES (
  gen_random_uuid(), NOW(), NOW(), $1
)

RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: DeleteAllUsers :exec
DELETE FROM users;

