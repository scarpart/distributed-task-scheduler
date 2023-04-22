-- name: CreateUser :one
INSERT INTO "Users" (
	username,
	password,
	email
) VALUES (
	$1, $2, $3
) RETURNING *;

-- name: GetUser :one
SELECT * FROM "Users"
WHERE user_id = $1 LIMIT $1;

-- name: GetAllUsers :many
SELECT * FROM "Users"
ORDER BY username
LIMIT $1
OFFSET $2;

-- name: DeleteUser :exec
DELETE FROM "Users"
WHERE user_id = $1;

-- name: UpdateUserPassword :one
UPDATE "Users"
SET password = $2
WHERE user_id = $1
RETURNING *; 

-- name: UpdateUserUsername :one
UPDATE "Users"
SET username = $2
WHERE user_id = $1
RETURNING *;
