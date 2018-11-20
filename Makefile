.PHONY: all
all: build
FORCE: ;

SHELL  := env BOOKMARK_ENV=$(BOOKMARK_ENV) $(SHELL)
BOOKMARK_ENV ?= dev

BIN_DIR = $(PWD)/bin

.PHONY: build

clean:
	rm -rf bin/*

dependencies:
	go mod download

build: dependencies build-api build-cmd

build-api: 
	go build -tags $(BOOKMARK_ENV) -o ./bin/api api/main.go

build-cmd:
	go build -tags $(BOOKMARK_ENV) -o ./bin/search cmd/main.go

linux-binaries:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags "$(BOOKMARK_ENV) netgo" -installsuffix netgo -o $(BIN_DIR)/api api/main.go
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags "$(BOOKMARK_ENV) netgo" -installsuffix netgo -o $(BIN_DIR)/search cmd/main.go

ci: dependencies test	

test:
	go test -tags testing ./...

fmt: ## gofmt and goimports all go files
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done