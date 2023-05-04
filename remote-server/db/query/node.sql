-- name: CreateNode :one
INSERT INTO "Nodes" (
	node_id,
	hostname,
	ip_addr,
	status
) VALUES (
	$1, $2, $3, $4
) RETURNING *;

-- name: GetNode :one
SELECT * FROM "Nodes"
WHERE node_id = $1 LIMIT 1;

-- name: GetAllNodes :many
SELECT * FROM "Nodes"
ORDER BY node_id
LIMIT $1
OFFSET $2; 

-- name: UpdateNode :one
UPDATE "Nodes"
SET hostname = $2, ip_addr = $3, status = $4
WHERE node_id = $1
RETURNING *;

