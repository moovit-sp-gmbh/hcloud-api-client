-- +goose Up
ALTER TABLE context
ADD COLUMN new_name TEXT;

-- +goose Down
ALTER TABLE context
DROP COLUMN new_name;
