.PHONY: run test lint tidy
 
 run:
	go run ./cmd/game

test:
	go test ./... -race -cover

lint:
	golangci-lint run ./..

format:
	gofmt -s -w .
	goimports -w .

tidy:
	go mod tidy