.PHONY: build test clean

build:
	go build -o biathlon ./cmd/main.go

test:
	go test ./internal/event ./internal/report ./internal/utils -v
