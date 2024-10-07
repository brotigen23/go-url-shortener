all: test run

run:
	go run cmd/shortener/main.go

test:
	go test ./internal/handlers