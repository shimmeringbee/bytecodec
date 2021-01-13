.DEFAULT_GOAL := all

build:
	go build

test:
	go test ./...

lint:
	go mod tidy
	golangci-lint run --fix -E goimports -E golint -E gochecknoinits -E whitespace -E gocyclo -E godox -E gocritic

install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint

clean:
	go clean -x

update:
	go get -u

all: clean lint test build