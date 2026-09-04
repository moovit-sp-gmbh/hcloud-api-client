-- +goose Up
ALTER TABLE context
ADD COLUMN oauth_app_client_id TEXT;

-- +goose Down
ALTER TABLE context
DROP COLUMN oauth_app_client_id;
