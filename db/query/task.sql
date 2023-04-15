-- name: CreateTask :one
INSERT INTO "Tasks" (
    user_id,
    task_name,
    task_description,
    status,
    priority,
    command
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetTask :one
SELECT * FROM "Tasks"
WHERE task_id = $1 LIMIT 1;

-- name: GetAllTasks :many
SELECT * FROM "Tasks"
ORDER BY task_id
LIMIT $1
OFFSET $2;

-- name: DeleteTask :exec
DELETE FROM "Tasks"
WHERE task_id = $1;