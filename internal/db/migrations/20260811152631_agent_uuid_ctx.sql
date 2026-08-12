-- +goose Up
ALTER TABLE context
ADD COLUMN agent_uuid TEXT;

-- +goose Down
ALTER TABLE context
DROP COLUMN agent_uuid;
