test:
	go test ./... && go mod tidy && go vet ./...

build:
	go build -o hcloud-cli ./main.go
	go build -o patcher ./cmd/patcher
	go build -o bin/codegen ./cmd/codegen

# Which endpoints have commands, which are left, and which the manifest claims
# but the clients no longer expose. Exits non-zero on the last.
coverage: bin/codegen
	./bin/codegen coverage --verify

missing: bin/codegen
	./bin/codegen coverage --missing

# Generate a command for an endpoint. Pass a service with SERVICE=cosmo.
new-command: bin/codegen
	./bin/codegen new $(if $(SERVICE),--client $(SERVICE),)

# Replay every recipe, e.g. after changing the emitter or regenerating a client.
regen-commands: bin/codegen
	./bin/codegen regen

bin/codegen: $(shell find cmd/codegen internal/codegen -name '*.go')
	go build -o bin/codegen ./cmd/codegen

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
	$(MAKE) generate+cosmo
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
