-- +goose Up
CREATE TABLE context (
	id INTEGER PRIMARY KEY,
	timestamp INTEGER NOT NULL DEFAULT (unixepoch()),
	jwt TEXT
);

-- +goose Down
DROP TABLE context;
