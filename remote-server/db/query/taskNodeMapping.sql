-- name: CreateTaskNodeMapping :one
INSERT INTO "TaskNodeMapping" (
	tnm_id,
	task_id,
	node_id
) VALUES (
	$1, $2, $3
) RETURNING *;

-- name: GetTaskNodeMapping :one
SELECT * FROM "TaskNodeMapping"
WHERE tnm_id = $1 LIMIT $1;

-- name: GetAllTaskNodeMappings :many
SELECT * FROM "TaskNodeMapping"
ORDER BY tnm_id
LIMIT $1 
OFFSET $2;

-- name: UpdateTaskNodeMapping :one
UPDATE "TaskNodeMapping"
SET task_id = $2, node_id = $3
WHERE tnm_id = $1
RETURNING *;


