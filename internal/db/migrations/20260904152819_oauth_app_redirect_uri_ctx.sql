-- +goose Up
ALTER TABLE context
ADD COLUMN oauth_app_redirect_uri TEXT;

-- +goose Down
ALTER TABLE context
DROP COLUMN oauth_app_redirect_uri;
