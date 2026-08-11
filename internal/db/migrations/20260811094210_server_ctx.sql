-- +goose Up
ALTER TABLE context
ADD COLUMN server TEXT;

-- +goose Down
ALTER TABLE context
DROP COLUMN server;
