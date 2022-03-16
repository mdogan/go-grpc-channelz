.PHONY: build test

build:
	go mod tidy
	go fmt ./...
	go vet ./...
	go build


test:
	go mod tidy
	go test ./...
