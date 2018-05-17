.PHONY: all
all: build
FORCE: ;

SHELL  := env BOOKMARK_ENV=$(BOOKMARK_ENV) $(SHELL)
BOOKMARK_ENV ?= dev

include config/$(BOOKMARK_ENV).env
export $(shell sed 's/=.*//' config/$(BOOKMARK_ENV).env)
BIN_DIR = $(PWD)/bin

.PHONY: build

clean:
	rm -rf bin/*

dependencies:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

build: dependencies build-api build-cmd

build-api: 
	go build -o ./bin/api api/main.go

build-cmd:
	go build -o ./bin/search cmd/main.go

linux-binaries:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags netgo -installsuffix netgo -o $(BIN_DIR)/api api/main.go
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags netgo -installsuffix netgo -o $(BIN_DIR)/search cmd/main.go

ci: dependencies test	

test:
	export BOOKMARK_ENV=$(BOOKMARK_ENV); go test  ./...

fmt: ## gofmt and goimports all go files
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done