-- +goose Up
ALTER TABLE context
ADD COLUMN agent_hello_token TEXT;

-- +goose Down
ALTER TABLE context
DROP COLUMN agent_hello_token;
