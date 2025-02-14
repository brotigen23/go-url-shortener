.PHONY: all
all: run


.PHONY: run
run:
	go run -ldflags "-X main.buildVersion=v0.0.1 -X 'main.buildDate=$(shell date)' -X main.buildCommit=asd123" cmd/shortener/main.go 


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
testCover:
	go test ./... -coverprofile tmp/coverage.out
	go tool cover -html=tmp/coverage.out -o tmp/cover.html

.PHONY: vet
vet:
	go run cmd/staticlint/multichecker.go ./...

build:
	go build -o cmd/shortener/shortener cmd/shortener/main.go

autoTests:
	./tmp/shortenertest \
	-source-path=. \
	-binary-path=cmd/shortener/shortener \
	-server-port=8080 \
	-file-storage-path=test/aliases.txt \
	-database-dsn='host=localhost port=5432 user=myuser password=1234 dbname=mydb sslmode=disable' \
	>tmp/.log

