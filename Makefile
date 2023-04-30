GOCMD  := go
GOFMT  := ${GOCMD} fmt
GOGET  := ${GOCMD} get
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

update:
	${GOGET} -u ./...

.PHONY: fmt lint test tidy update
