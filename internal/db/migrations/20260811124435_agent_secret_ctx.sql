-- +goose Up
ALTER TABLE context
ADD COLUMN agent_secret TEXT;

-- +goose Down
ALTER TABLE context
DROP COLUMN agent_secret;
