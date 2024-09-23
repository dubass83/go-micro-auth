-- name: Insert :one
INSERT INTO users (
  email, first_name, last_name, password, active
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetOne :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetAll :many
SELECT * FROM users
ORDER by last_name
LIMIT $1
OFFSET $2;

-- name: Update :one
UPDATE users 
SET 
  first_name = COALESCE(sqlc.narg('first_name'), first_name),
  last_name = COALESCE(sqlc.narg('last_name'), last_name),
  email = COALESCE(sqlc.narg('email'), email),
  active = COALESCE(sqlc.narg('active'), active),
  updated_at = COALESCE(sqlc.narg('updated_at'), updated_at)
WHERE 
  id = sqlc.arg('id')
RETURNING *;

-- name: Delete :exec
DELETE FROM users
WHERE id = $1; 

-- name: DeleteByID :exec
DELETE FROM users
WHERE id = $1;

