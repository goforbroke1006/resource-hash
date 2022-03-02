all: dep build test lint

.PHONY: dep build test lint

dep:
	go mod download

build:
	go build ./

test:
	go test ./... -cover

lint:
	golangci-lint run
