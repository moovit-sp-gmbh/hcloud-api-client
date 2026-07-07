-- name: CreateContext :exec
INSERT INTO context (timestamp, JWT, org_name) VALUES (?, ?, ?);

-- name: GetContext :one
SELECT * FROM context
WHERE id=?;

-- name: GetLatestContext :one
SELECT * FROM context
ORDER BY timestamp DESC
LIMIT 1;

-- name: List :many
SELECT * FROM context
ORDER BY timestamp DESC
LIMIT ?
OFFSET ?;

-- name: SetTimestamp :exec
UPDATE context
SET timestamp = ?
WHERE ID = ?;
