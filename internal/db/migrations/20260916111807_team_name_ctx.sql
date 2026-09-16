-- +goose Up
ALTER TABLE context
ADD COLUMN team_name TEXT;

-- +goose Down
ALTER TABLE context
DROP COLUMN team_name;
