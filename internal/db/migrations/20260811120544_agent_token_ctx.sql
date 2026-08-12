-- +goose Up
ALTER TABLE context
ADD COLUMN agent_token TEXT;

-- +goose Down
ALTER TABLE context
DROP COLUMN agent_token;
