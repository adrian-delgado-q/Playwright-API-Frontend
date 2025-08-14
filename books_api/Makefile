build:
	go build -o bin/books_api main.go

run:
	go run main.go

test:
	go test -v

test-coverage:
	go test -v -cover

deps:
	go mod tidy

clean:
	rm -f bin/books_api books.db

dev:
	go run main.go

.PHONY: build run test test-coverage deps clean dev
