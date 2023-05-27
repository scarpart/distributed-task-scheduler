-- name: CreateLog :one
INSERT INTO "Logs" (
	task_id,
	message
) VALUES (
	$1, $2
) RETURNING *;

-- name: UpdateLog :one
UPDATE "Logs" 
SET message = $2
WHERE log_id = $1
RETURNING *;
