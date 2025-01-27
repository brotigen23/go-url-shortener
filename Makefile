.PHONY: all run test client

all: run

run:
	go run cmd/shortener/main.go

test:
	go test ./... -v -cover 

mock:
	~/go/bin/mockgen -destination=internal/mock/mockRepository.go -package=mock -source=internal/repository/repository.go


