.PHONY: all
all: run


.PHONY: run
run:
	go run cmd/shortener/main.go > ./errors.log

.PHONY: client
client:
	go run cmd/client/main.go

.PHONY: test
test:
	go test ./... -v -cover -count=1

.PHONY: mock
mock:
	~/go/bin/mockgen -destination=internal/mock/mockRepository.go -package=mock -source=internal/repository/repository.go

.PHONY: doc
doc:
	~/go/bin/godoc -http=:8080

.PHONY: testCover
testCover: coverage.out
	go tool cover -html=coverage.out -o tmp/cover.html

coverage.out:
	go test ./... -coverprofile coverage.out
