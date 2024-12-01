.PHONY: all run test client

all: run

run:
	go run cmd/shortener/main.go -d "host=localhost port=5432 user=myuser password=1234 dbname=mydb sslmode=disable"
r:
	go run cmd/shortener/main.go

test:
	go test ./internal/handlers -v -cover 

testM:
	cd internal/repositories && go test inMemoryRepository_test.go

client:
	go run cmd/client/main.go

testPQ:
	cd internal/repositories && go test ./... -v -count=1
testS:
	cd internal/services && go test ./... -v -count=1