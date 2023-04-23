GOCMD  := go
GOFMT  := ${GOCMD} fmt
GOMOD  := ${GOCMD} mod
GOTEST := ${GOCMD} test

fmt:
	${GOFMT} ./...

lint:
	golangci-lint run ./...

test:
	${GOTEST} ./...

tidy:
	${GOMOD} tidy

.PHONY: fmt lint test tidy
