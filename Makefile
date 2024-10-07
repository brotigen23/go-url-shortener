all: test run

run:
	go run cmd/shortener/server/main.go

test:
	go test ./internal/handlers

client:
	go run cmd/shortener/client/main.go