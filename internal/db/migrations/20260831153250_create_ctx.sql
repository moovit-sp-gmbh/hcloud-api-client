-- +goose Up
ALTER TABLE context
ADD COLUMN high5_space_id TEXT;
ALTER TABLE context
ADD COLUMN high5_space_name TEXT;

-- +goose Down
ALTER TABLE context
DROP COLUMN high5_space_id;
ALTER TABLE context
DROP COLUMN high5_space_name;
