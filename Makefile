.PHONY: build

test:
	go test ./...

run:
	go run -race main.go

build+linux:
	GOOS=linux GOARCH=amd64 go build -o bin/linux/hcloud main.go

build+linux+arm:
	GOOS=linux GOARCH=arm go build -o bin/linux/arm/hcloud main.go

build+darwin:
	GOOS=darwin GOARCH=amd64 go build -o bin/darwin/hcloud main.go

build+windows:
	GOOS=windows GOARCH=amd64 go build -o bin/windows/hcloud.exe main.go
