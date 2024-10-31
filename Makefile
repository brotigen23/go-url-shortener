.PHONY: all run test client postgresRun

all: run

run:
	go run cmd/shortener/main.go

test:
	go test ./internal/handlers -v -cover

client:
	go run cmd/client/main.go

postgresRun:
	psql -h localhost -U myuser -d mydb

t:
	curl --header "Content-Type: application/json" \
  	--request POST \
  	--data '{"url":"ya.ru"}' \
  	localhost:8080/api/shorten