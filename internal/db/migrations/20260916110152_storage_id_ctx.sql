-- +goose Up
ALTER TABLE context
ADD COLUMN storage_id TEXT;

-- +goose Down
ALTER TABLE context
DROP COLUMN storage_id;
