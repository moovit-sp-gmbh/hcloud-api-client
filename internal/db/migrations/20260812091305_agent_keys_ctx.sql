-- +goose Up
ALTER TABLE context
ADD COLUMN agent_private_key TEXT;
ALTER TABLE context
ADD COLUMN agent_public_key TEXT;

-- +goose Down
ALTER TABLE context
DROP COLUMN agent_private_key;
ALTER TABLE context
DROP COLUMN agent_public_key;
