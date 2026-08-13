test:
	go test ./... && go mod tidy && go vet ./...

build:
	go build -o hcloud-cli ./main.go
	go build -o patcher ./cmd/patcher

format:
	go fmt ./...

sqlc:
	sqlc generate

migrate+create:
	read -p "Migration name: "; goose sqlite3 db.sqlite -dir internal/db/migrations create $$REPLY sql

migrate+down:
	goose sqlite3 db.sqlite down -dir internal/db/migrations

migrate+up:
	goose sqlite3 db.sqlite up -dir internal/db/migrations

generate:
	$(MAKE) generate+idp
	$(MAKE) generate+high5
	$(MAKE) generate+mothership
#	$(MAKE) generate+cosmo
	$(MAKE) generate+dali

generate+idp:
	./scripts/generate_idp.sh

generate+high5:
	./scripts/generate_high5.sh

generate+mothership:
	./scripts/generate_mothership.sh

generate+cosmo:
	./scripts/generate_cosmo.sh

generate+dali:
	./scripts/generate_dali.sh
