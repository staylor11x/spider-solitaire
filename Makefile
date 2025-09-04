.PHONY: run test lint tidy
 
 run:
	go run ./cmd/gane

test:
	go test ./... -race -cover

lint:
	golangci-lint run ./..

tidy:
	go mod tidy