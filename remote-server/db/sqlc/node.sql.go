// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: node.sql

package db

import (
	"context"
)

const createNode = `-- name: CreateNode :one
INSERT INTO "Nodes" (
	hostname,
	ip_addr,
	status
) VALUES (
	$1, $2, $3
) RETURNING node_id, hostname, ip_addr, status, created_at, updated_at
`

type CreateNodeParams struct {
	Hostname string `json:"hostname"`
	IpAddr   string `json:"ip_addr"`
	Status   int32  `json:"status"`
}

func (q *Queries) CreateNode(ctx context.Context, arg CreateNodeParams) (Node, error) {
	row := q.db.QueryRowContext(ctx, createNode, arg.Hostname, arg.IpAddr, arg.Status)
	var i Node
	err := row.Scan(
		&i.NodeID,
		&i.Hostname,
		&i.IpAddr,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getAllNodes = `-- name: GetAllNodes :many
SELECT node_id, hostname, ip_addr, status, created_at, updated_at FROM "Nodes"
ORDER BY node_id
LIMIT $1
OFFSET $2
`

type GetAllNodesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetAllNodes(ctx context.Context, arg GetAllNodesParams) ([]Node, error) {
	rows, err := q.db.QueryContext(ctx, getAllNodes, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Node
	for rows.Next() {
		var i Node
		if err := rows.Scan(
			&i.NodeID,
			&i.Hostname,
			&i.IpAddr,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const getNode = `-- name: GetNode :one
SELECT node_id, hostname, ip_addr, status, created_at, updated_at FROM "Nodes"
WHERE node_id = $1 LIMIT 1
`

func (q *Queries) GetNode(ctx context.Context, nodeID int64) (Node, error) {
	row := q.db.QueryRowContext(ctx, getNode, nodeID)
	var i Node
	err := row.Scan(
		&i.NodeID,
		&i.Hostname,
		&i.IpAddr,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateNode = `-- name: UpdateNode :one
UPDATE "Nodes"
SET hostname = $2, ip_addr = $3, status = $4
WHERE node_id = $1
RETURNING node_id, hostname, ip_addr, status, created_at, updated_at
`

type UpdateNodeParams struct {
	NodeID   int64  `json:"node_id"`
	Hostname string `json:"hostname"`
	IpAddr   string `json:"ip_addr"`
	Status   int32  `json:"status"`
}

func (q *Queries) UpdateNode(ctx context.Context, arg UpdateNodeParams) (Node, error) {
	row := q.db.QueryRowContext(ctx, updateNode,
		arg.NodeID,
		arg.Hostname,
		arg.IpAddr,
		arg.Status,
	)
	var i Node
	err := row.Scan(
		&i.NodeID,
		&i.Hostname,
		&i.IpAddr,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
