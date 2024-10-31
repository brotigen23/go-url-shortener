.PHONY: all run test client postgresRun

all: run

run:
	go run cmd/shortener/main.go

test:
	go test ./internal/handlers

client:
	go run cmd/client/main.go

postgresRun:
	psql -h localhost -U myuser -d mydb