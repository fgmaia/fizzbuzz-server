.PHONY: mocks
mocks:
	go install github.com/vektra/mockery/v2@latest
	mockery --all

.PHONY: test
test:
	go test ./... -v

.PHONY: build
build:
	go build -o bin/fizzbuzz-server ./cmd/server

.PHONY: run
run:
	go run ./cmd/server