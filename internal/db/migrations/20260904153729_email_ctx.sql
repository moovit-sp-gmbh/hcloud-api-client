-- +goose Up
ALTER TABLE context
ADD COLUMN email TEXT;

-- +goose Down
ALTER TABLE context
DROP COLUMN email;
