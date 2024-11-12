.PHONY: all run test client postgresRun

all: run

run:
	go run cmd/shortener/main.go

test:
	go test ./internal/handlers -v -cover 

client:
	go run cmd/client/main.go

setENV:
	DATABASE_DSN="host=localhost port=5432 user=myuser password=1234 dbname=mydb sslmode=disable"

t:
	curl -i --header "Content-Type: application/json" \
  	--request POST \
  	--data '{"url":"ya.ru"}' \
  	localhost:8080/api/shorten
	curl -i --header "Content-Type: application/json" \
  	--request POST \
  	--data '{"url":"yandex.ru"}' \
  	localhost:8080/api/shorten
	curl -i --header "Content-Type: application/json" \
  	--request POST \
  	--data '{"url":"google.com"}' \
  	localhost:8080/api/shorten
	curl -i --header "Content-Type: application/json" \
  	--request POST \
  	--data '{"url":"metanit.net"}' \
  	localhost:8080/api/shorten
	curl -i --header "Content-Type: application/json" \
  	--request POST \
  	--data '{"url":"habr.com"}' \
  	localhost:8080/api/shorten