.PHONY: all run test client

all: run

run:
	go run cmd/shortener/main.go -d "host=localhost port=5432 user=myuser password=1234 dbname=mydb sslmode=disable"

test:
	go test ./internal/handlers -v -cover 

client:
	go run cmd/client/main.go
