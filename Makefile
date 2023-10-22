.PHONY: build test

build:
	go build -o ratingd ./cmd/ratingd

test:
	go test ./...

cover:
	go test -coverprofile=cover.out ./...
