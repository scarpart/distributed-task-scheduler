-- name: CreateLog :one
INSERT INTO "Logs" (
	log_id,
	task_id,
	message
) VALUES (
	$1, $2, $3
) RETURNING *;

-- name: UpdateLog :one
UPDATE "Logs" 
SET message = $2
WHERE log_id = $1
RETURNING *;
