-- name: CreateUser :one
INSERT INTO users (
  username,
  email,
  hashed_password
) VALUES (
  $1,
  $2,
  $3
) RETURNING *;

-- name: GetUserById :one
SELECT *
FROM users
WHERE id = $1;