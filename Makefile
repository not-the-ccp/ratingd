.PHONY: build test

build:
	go build -o ratingd ./cmd/ratingd

test:
	go test ./...
