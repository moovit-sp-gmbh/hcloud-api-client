-- +goose Up
ALTER TABLE context
ADD COLUMN org_name TEXT;

-- +goose Down
ALTER TABLE context
DROP COLUMN org_name;
