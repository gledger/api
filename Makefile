run:
	go run cmd/api/*.go

test:
	go test -v .
	go test -v ./cmd/api
	go test -v ./db

lint:
	golint .
	golint ./cmd/api
	golint ./db

.PHONY: run lint test
