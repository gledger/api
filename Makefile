run:
	go run cmd/api/*.go

test: vet
	@go test -v .
	@go test -v ./cmd/api
	@go test -v ./db

vet:
	@go vet .
	@go vet ./cmd/api
	@go vet ./db

lint:
	golint .
	golint ./cmd/api
	golint ./db

.PHONY: run lint test
