format:
	go fmt ./...

migrate+create:
	read -p "Migration name: "; goose sqlite3 db.sqlite -dir internal/db/migrations create $$REPLY sql

migrate+down:
	goose sqlite3 db.sqlite down -dir internal/db/migrations

migrate+up:
	goose sqlite3 db.sqlite up -dir internal/db/migrations
