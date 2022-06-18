all: lint test build
.PHONY: all

build:
	go build ./cmd/gofixit
.PHONY: build

test:
	go test ./... -coverprofile cover.out
	go tool cover -func=cover.out
.PHONY: test

lint:
	go vet ./...
.PHONY: lint

tidy:
	go mod tidy
.PHONY: tidy

clean:
	rm -rf cover.out
.PHONY: clean
