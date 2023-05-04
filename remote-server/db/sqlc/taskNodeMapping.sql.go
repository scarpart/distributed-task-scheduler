// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: taskNodeMapping.sql

package db

import (
	"context"
)

const createTaskNodeMapping = `-- name: CreateTaskNodeMapping :one
INSERT INTO "TaskNodeMapping" (
	tnm_id,
	task_id,
	node_id
) VALUES (
	$1, $2, $3
) RETURNING tnm_id, task_id, node_id, created_at
`

type CreateTaskNodeMappingParams struct {
	TnmID  int64 `json:"tnm_id"`
	TaskID int64 `json:"task_id"`
	NodeID int64 `json:"node_id"`
}

func (q *Queries) CreateTaskNodeMapping(ctx context.Context, arg CreateTaskNodeMappingParams) (TaskNodeMapping, error) {
	row := q.db.QueryRowContext(ctx, createTaskNodeMapping, arg.TnmID, arg.TaskID, arg.NodeID)
	var i TaskNodeMapping
	err := row.Scan(
		&i.TnmID,
		&i.TaskID,
		&i.NodeID,
		&i.CreatedAt,
	)
	return i, err
}

const getAllTaskNodeMappings = `-- name: GetAllTaskNodeMappings :many
SELECT tnm_id, task_id, node_id, created_at FROM "TaskNodeMapping"
ORDER BY tnm_id
LIMIT $1 
OFFSET $2
`

type GetAllTaskNodeMappingsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetAllTaskNodeMappings(ctx context.Context, arg GetAllTaskNodeMappingsParams) ([]TaskNodeMapping, error) {
	rows, err := q.db.QueryContext(ctx, getAllTaskNodeMappings, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []TaskNodeMapping
	for rows.Next() {
		var i TaskNodeMapping
		if err := rows.Scan(
			&i.TnmID,
			&i.TaskID,
			&i.NodeID,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTaskNodeMapping = `-- name: GetTaskNodeMapping :one
SELECT tnm_id, task_id, node_id, created_at FROM "TaskNodeMapping"
WHERE tnm_id = $1 LIMIT $1
`

func (q *Queries) GetTaskNodeMapping(ctx context.Context, limit int32) (TaskNodeMapping, error) {
	row := q.db.QueryRowContext(ctx, getTaskNodeMapping, limit)
	var i TaskNodeMapping
	err := row.Scan(
		&i.TnmID,
		&i.TaskID,
		&i.NodeID,
		&i.CreatedAt,
	)
	return i, err
}

const updateTaskNodeMapping = `-- name: UpdateTaskNodeMapping :one
UPDATE "TaskNodeMapping"
SET task_id = $2, node_id = $3
WHERE tnm_id = $1
RETURNING tnm_id, task_id, node_id, created_at
`

type UpdateTaskNodeMappingParams struct {
	TnmID  int64 `json:"tnm_id"`
	TaskID int64 `json:"task_id"`
	NodeID int64 `json:"node_id"`
}

func (q *Queries) UpdateTaskNodeMapping(ctx context.Context, arg UpdateTaskNodeMappingParams) (TaskNodeMapping, error) {
	row := q.db.QueryRowContext(ctx, updateTaskNodeMapping, arg.TnmID, arg.TaskID, arg.NodeID)
	var i TaskNodeMapping
	err := row.Scan(
		&i.TnmID,
		&i.TaskID,
		&i.NodeID,
		&i.CreatedAt,
	)
	return i, err
}