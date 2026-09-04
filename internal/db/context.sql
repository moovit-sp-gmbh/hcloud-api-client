-- name: CreateContext :exec
INSERT INTO context (
        timestamp,
        JWT,
        org_name,
        new_name,
        server,
        agent_token,
        agent_secret,
        agent_uuid,
        agent_private_key,
        agent_public_key,
        agent_hello_token,
        high5_space_id,
        high5_space_name,
        oauth_app_client_id,
        oauth_app_redirect_uri
) VALUES (
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?
);

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
