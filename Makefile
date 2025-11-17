.PHONY: build test

build:
	go build -o ratingd ./cmd/ratingd

test:
	go test ./...

cover:
	go test -coverprofile=cover.out ./...

tidy:
	go mod tidy

migrate:
	goose -dir migrations postgres "$(DATABASE_URL)" up

migrate:
	goose -dir migrations postgres "$(DATABASE_URL)" up

migrate:
	goose -dir migrations postgres "$(DATABASE_URL)" up

migrate:
	goose -dir migrations postgres "$(DATABASE_URL)" up

cover:
	go test -coverprofile=cover.out ./...

tidy:
	go mod tidy
