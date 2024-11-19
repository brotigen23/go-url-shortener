.PHONY: all run test client postgresRun migrate

all: run

run:
	go run cmd/shortener/main.go

test:
	go test ./internal/handlers -v -cover 

client:
	go run cmd/client/main.go

migrate:
	~/go/bin/goose -dir internal/db/migrations postgres  "$(DATABASE_DSN)" up